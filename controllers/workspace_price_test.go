package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"mobile-backend-go/constants"
	"mobile-backend-go/database"
	"mobile-backend-go/models"
)

type workspacePriceFixture struct {
	User              models.User
	PersonalWorkspace models.Workspace
	SecondWorkspace   models.Workspace
	Ingredient        models.Ingredient
	Recipe            models.Recipe
	SecondRecipe      models.Recipe
}

func setupWorkspacePriceTest(t *testing.T) workspacePriceFixture {
	t.Helper()

	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	t.Cleanup(func() {
		sqlDB, err := db.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	})

	if err := db.AutoMigrate(
		&models.User{},
		&models.Workspace{},
		&models.WorkspaceMember{},
		&models.Ingredient{},
		&models.Price{},
		&models.Recipe{},
		&models.RecipeIngredient{},
	); err != nil {
		t.Fatalf("migrate test database: %v", err)
	}
	database.DB = db

	user := models.User{Username: "phase2a-test", Password: "hashed"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}

	personalUserID := user.ID
	personalWorkspace := models.Workspace{
		Name:           "Personal workspace",
		Slug:           "personal-test",
		PersonalUserID: &personalUserID,
	}
	secondWorkspace := models.Workspace{
		Name: "Shared kitchen",
		Slug: "shared-kitchen",
	}
	if err := db.Create(&personalWorkspace).Error; err != nil {
		t.Fatalf("create personal workspace: %v", err)
	}
	if err := db.Create(&secondWorkspace).Error; err != nil {
		t.Fatalf("create second workspace: %v", err)
	}

	members := []models.WorkspaceMember{
		{WorkspaceID: personalWorkspace.ID, UserID: user.ID, Role: constants.WorkspaceRoleOwner},
		{WorkspaceID: secondWorkspace.ID, UserID: user.ID, Role: constants.WorkspaceRoleOwner},
	}
	if err := db.Create(&members).Error; err != nil {
		t.Fatalf("create workspace memberships: %v", err)
	}

	ingredient := models.Ingredient{Name: "Test flour", Type: "base"}
	if err := db.Create(&ingredient).Error; err != nil {
		t.Fatalf("create ingredient: %v", err)
	}

	recipe := models.Recipe{Name: "Test recipe", UserID: user.ID, WorkspaceID: &personalWorkspace.ID}
	if err := db.Create(&recipe).Error; err != nil {
		t.Fatalf("create recipe: %v", err)
	}
	secondRecipe := models.Recipe{Name: "Second workspace recipe", UserID: user.ID, WorkspaceID: &secondWorkspace.ID}
	if err := db.Create(&secondRecipe).Error; err != nil {
		t.Fatalf("create second recipe: %v", err)
	}
	recipeIngredient := models.RecipeIngredient{
		RecipeID:     recipe.ID,
		IngredientID: ingredient.ID,
		Quantity:     "1000",
		Unit:         "g",
	}
	if err := db.Create(&recipeIngredient).Error; err != nil {
		t.Fatalf("create recipe ingredient: %v", err)
	}
	secondRecipeIngredient := models.RecipeIngredient{
		RecipeID:     secondRecipe.ID,
		IngredientID: ingredient.ID,
		Quantity:     "1000",
		Unit:         "g",
	}
	if err := db.Create(&secondRecipeIngredient).Error; err != nil {
		t.Fatalf("create second recipe ingredient: %v", err)
	}

	return workspacePriceFixture{
		User:              user,
		PersonalWorkspace: personalWorkspace,
		SecondWorkspace:   secondWorkspace,
		Ingredient:        ingredient,
		Recipe:            recipe,
		SecondRecipe:      secondRecipe,
	}
}

func runWorkspaceRequest(
	userID uint,
	workspaceID uint,
	handler gin.HandlerFunc,
	method string,
	routePattern string,
	requestTarget string,
) *httptest.ResponseRecorder {
	router := gin.New()
	router.Handle(method, routePattern, func(c *gin.Context) {
		c.Set("userID", userID)
		c.Set("workspaceID", workspaceID)
		handler(c)
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(method, requestTarget, nil)
	router.ServeHTTP(recorder, request)
	return recorder
}

func TestGetPricesFiltersByWorkspace(t *testing.T) {
	fixture := setupWorkspacePriceTest(t)

	createPrice(t, fixture.User.ID, fixture.PersonalWorkspace.ID, fixture.Ingredient.ID, 10)
	createPrice(t, fixture.User.ID, fixture.SecondWorkspace.ID, fixture.Ingredient.ID, 20)

	personalResponse := runWorkspaceRequest(
		fixture.User.ID,
		fixture.PersonalWorkspace.ID,
		GetPrices,
		http.MethodGet,
		"/prices",
		"/prices?ingredient_id="+uintToString(fixture.Ingredient.ID),
	)
	if personalResponse.Code != http.StatusOK {
		t.Fatalf("personal workspace status = %d body = %s", personalResponse.Code, personalResponse.Body.String())
	}
	personalPrices := decodePrices(t, personalResponse)
	if len(personalPrices) != 1 {
		t.Fatalf("personal workspace price count = %d, want 1", len(personalPrices))
	}
	if personalPrices[0].Price != 10 {
		t.Fatalf("personal workspace price = %v, want 10", personalPrices[0].Price)
	}
	if personalPrices[0].WorkspaceID == nil || *personalPrices[0].WorkspaceID != fixture.PersonalWorkspace.ID {
		t.Fatalf("personal workspace id = %v, want %d", personalPrices[0].WorkspaceID, fixture.PersonalWorkspace.ID)
	}

	secondResponse := runWorkspaceRequest(
		fixture.User.ID,
		fixture.SecondWorkspace.ID,
		GetPrices,
		http.MethodGet,
		"/prices",
		"/prices?ingredient_id="+uintToString(fixture.Ingredient.ID),
	)
	if secondResponse.Code != http.StatusOK {
		t.Fatalf("second workspace status = %d body = %s", secondResponse.Code, secondResponse.Body.String())
	}
	secondPrices := decodePrices(t, secondResponse)
	if len(secondPrices) != 1 {
		t.Fatalf("second workspace price count = %d, want 1", len(secondPrices))
	}
	if secondPrices[0].Price != 20 {
		t.Fatalf("second workspace price = %v, want 20", secondPrices[0].Price)
	}
	if secondPrices[0].WorkspaceID == nil || *secondPrices[0].WorkspaceID != fixture.SecondWorkspace.ID {
		t.Fatalf("second workspace id = %v, want %d", secondPrices[0].WorkspaceID, fixture.SecondWorkspace.ID)
	}
}

func TestGetRecipeCostUsesActiveWorkspacePrice(t *testing.T) {
	fixture := setupWorkspacePriceTest(t)

	createPrice(t, fixture.User.ID, fixture.PersonalWorkspace.ID, fixture.Ingredient.ID, 10)
	createPrice(t, fixture.User.ID, fixture.SecondWorkspace.ID, fixture.Ingredient.ID, 20)

	personalRecipe := getRecipeForWorkspace(t, fixture, fixture.PersonalWorkspace.ID)
	if personalRecipe.TotalCost != 10 {
		t.Fatalf("personal workspace recipe total = %v, want 10", personalRecipe.TotalCost)
	}
	assertAttachedLatestPrice(t, personalRecipe, 10, fixture.PersonalWorkspace.ID)

	secondRecipe := getRecipeForWorkspace(t, fixture, fixture.SecondWorkspace.ID, fixture.SecondRecipe.ID)
	if secondRecipe.TotalCost != 20 {
		t.Fatalf("second workspace recipe total = %v, want 20", secondRecipe.TotalCost)
	}
	assertAttachedLatestPrice(t, secondRecipe, 20, fixture.SecondWorkspace.ID)
}

func TestGetRecipeCostIsZeroWhenWorkspaceHasNoPrice(t *testing.T) {
	fixture := setupWorkspacePriceTest(t)

	createPrice(t, fixture.User.ID, fixture.PersonalWorkspace.ID, fixture.Ingredient.ID, 10)

	recipe := getRecipeForWorkspace(t, fixture, fixture.SecondWorkspace.ID, fixture.SecondRecipe.ID)
	if recipe.TotalCost != 0 {
		t.Fatalf("empty workspace recipe total = %v, want 0", recipe.TotalCost)
	}
	if len(recipe.RecipeIngredients) != 1 {
		t.Fatalf("recipe ingredient count = %d, want 1", len(recipe.RecipeIngredients))
	}
	if len(recipe.RecipeIngredients[0].Ingredient.Prices) != 0 {
		t.Fatalf("empty workspace attached price count = %d, want 0", len(recipe.RecipeIngredients[0].Ingredient.Prices))
	}
}

func createPrice(t *testing.T, userID uint, workspaceID uint, ingredientID uint, value float64) {
	t.Helper()

	price := models.Price{
		IngredientID: ingredientID,
		Price:        value,
		Quantity:     1,
		Unit:         "kg",
		Date:         time.Now(),
		UserID:       userID,
		WorkspaceID:  &workspaceID,
	}
	if err := database.DB.Create(&price).Error; err != nil {
		t.Fatalf("create price: %v", err)
	}
}

func getRecipeForWorkspace(t *testing.T, fixture workspacePriceFixture, workspaceID uint, recipeID ...uint) models.Recipe {
	t.Helper()

	targetRecipeID := fixture.Recipe.ID
	if len(recipeID) > 0 {
		targetRecipeID = recipeID[0]
	}

	response := runWorkspaceRequest(
		fixture.User.ID,
		workspaceID,
		GetRecipe,
		http.MethodGet,
		"/recipes/:id",
		"/recipes/"+uintToString(targetRecipeID),
	)
	if response.Code != http.StatusOK {
		t.Fatalf("get recipe status = %d body = %s", response.Code, response.Body.String())
	}

	var recipe models.Recipe
	if err := json.Unmarshal(response.Body.Bytes(), &recipe); err != nil {
		t.Fatalf("decode recipe response: %v", err)
	}
	return recipe
}

func decodePrices(t *testing.T, response *httptest.ResponseRecorder) []models.Price {
	t.Helper()

	var prices []models.Price
	if err := json.Unmarshal(response.Body.Bytes(), &prices); err != nil {
		t.Fatalf("decode price response: %v", err)
	}
	return prices
}

func assertAttachedLatestPrice(t *testing.T, recipe models.Recipe, wantPrice float64, wantWorkspaceID uint) {
	t.Helper()

	if len(recipe.RecipeIngredients) != 1 {
		t.Fatalf("recipe ingredient count = %d, want 1", len(recipe.RecipeIngredients))
	}
	prices := recipe.RecipeIngredients[0].Ingredient.Prices
	if len(prices) != 1 {
		t.Fatalf("attached latest price count = %d, want 1", len(prices))
	}
	if prices[0].Price != wantPrice {
		t.Fatalf("attached latest price = %v, want %v", prices[0].Price, wantPrice)
	}
	if prices[0].WorkspaceID == nil || *prices[0].WorkspaceID != wantWorkspaceID {
		t.Fatalf("attached latest price workspace = %v, want %d", prices[0].WorkspaceID, wantWorkspaceID)
	}
}

func uintToString(value uint) string {
	return strconv.FormatUint(uint64(value), 10)
}
