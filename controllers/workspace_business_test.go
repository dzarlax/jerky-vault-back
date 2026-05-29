package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"mobile-backend-go/constants"
	"mobile-backend-go/database"
	"mobile-backend-go/models"
)

type workspaceBusinessFixture struct {
	User              models.User
	PersonalWorkspace models.Workspace
	SecondWorkspace   models.Workspace
	Ingredient        models.Ingredient
	PersonalRecipe    models.Recipe
	SecondRecipe      models.Recipe
	PersonalPackage   models.Package
	SecondPackage     models.Package
	PersonalProduct   models.Product
	SecondProduct     models.Product
	PersonalClient    models.Client
	SecondClient      models.Client
	PersonalOrder     models.Order
	SecondOrder       models.Order
	PersonalSession   models.CookingSession
	SecondSession     models.CookingSession
}

func setupWorkspaceBusinessTest(t *testing.T) workspaceBusinessFixture {
	t.Helper()

	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
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
		&models.Package{},
		&models.Product{},
		&models.ProductOption{},
		&models.Client{},
		&models.Order{},
		&models.OrderItem{},
		&models.CookingSession{},
		&models.CookingSessionIngredient{},
	); err != nil {
		t.Fatalf("migrate test database: %v", err)
	}
	database.DB = db

	user := models.User{Username: "phase2b-test", Password: "hashed"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}

	personalUserID := user.ID
	personalWorkspace := models.Workspace{Name: "Personal", Slug: "personal-business", PersonalUserID: &personalUserID}
	secondWorkspace := models.Workspace{Name: "Second", Slug: "second-business"}
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
		t.Fatalf("create memberships: %v", err)
	}

	ingredient := models.Ingredient{Name: "Pepper", Type: "spice"}
	if err := db.Create(&ingredient).Error; err != nil {
		t.Fatalf("create ingredient: %v", err)
	}

	personalRecipe := models.Recipe{Name: "Personal recipe", UserID: user.ID, WorkspaceID: &personalWorkspace.ID}
	secondRecipe := models.Recipe{Name: "Second recipe", UserID: user.ID, WorkspaceID: &secondWorkspace.ID}
	personalPackage := models.Package{Name: "Personal pack", UserID: user.ID, WorkspaceID: &personalWorkspace.ID}
	secondPackage := models.Package{Name: "Second pack", UserID: user.ID, WorkspaceID: &secondWorkspace.ID}
	personalClient := models.Client{Name: "Personal", Surname: "Client", UserID: user.ID, WorkspaceID: &personalWorkspace.ID}
	secondClient := models.Client{Name: "Second", Surname: "Client", UserID: user.ID, WorkspaceID: &secondWorkspace.ID}
	if err := db.Create(&personalRecipe).Error; err != nil {
		t.Fatalf("create personal recipe: %v", err)
	}
	if err := db.Create(&secondRecipe).Error; err != nil {
		t.Fatalf("create second recipe: %v", err)
	}
	if err := db.Create(&personalPackage).Error; err != nil {
		t.Fatalf("create personal package: %v", err)
	}
	if err := db.Create(&secondPackage).Error; err != nil {
		t.Fatalf("create second package: %v", err)
	}
	if err := db.Create(&personalClient).Error; err != nil {
		t.Fatalf("create personal client: %v", err)
	}
	if err := db.Create(&secondClient).Error; err != nil {
		t.Fatalf("create second client: %v", err)
	}

	personalProduct := models.Product{Name: "Personal product", Price: 12, Cost: 5, UserID: user.ID, WorkspaceID: &personalWorkspace.ID, PackageID: personalPackage.ID}
	secondProduct := models.Product{Name: "Second product", Price: 21, Cost: 8, UserID: user.ID, WorkspaceID: &secondWorkspace.ID, PackageID: secondPackage.ID}
	if err := db.Create(&personalProduct).Error; err != nil {
		t.Fatalf("create personal product: %v", err)
	}
	if err := db.Create(&secondProduct).Error; err != nil {
		t.Fatalf("create second product: %v", err)
	}

	personalOrder := models.Order{ClientID: personalClient.ID, Date: time.Now(), Status: constants.OrderStatusFinished, UserID: user.ID, WorkspaceID: &personalWorkspace.ID}
	secondOrder := models.Order{ClientID: secondClient.ID, Date: time.Now().Add(time.Hour), Status: constants.OrderStatusNew, UserID: user.ID, WorkspaceID: &secondWorkspace.ID}
	if err := db.Create(&personalOrder).Error; err != nil {
		t.Fatalf("create personal order: %v", err)
	}
	if err := db.Create(&secondOrder).Error; err != nil {
		t.Fatalf("create second order: %v", err)
	}
	items := []models.OrderItem{
		{OrderID: personalOrder.ID, ProductID: personalProduct.ID, Quantity: 2, Price: 12, Cost_price: 5},
		{OrderID: secondOrder.ID, ProductID: secondProduct.ID, Quantity: 3, Price: 21, Cost_price: 8},
	}
	if err := db.Create(&items).Error; err != nil {
		t.Fatalf("create order items: %v", err)
	}

	personalSession := models.CookingSession{RecipeID: personalRecipe.ID, Date: time.Now(), Yield: "1 kg", UserID: user.ID, WorkspaceID: &personalWorkspace.ID}
	secondSession := models.CookingSession{RecipeID: secondRecipe.ID, Date: time.Now(), Yield: "2 kg", UserID: user.ID, WorkspaceID: &secondWorkspace.ID}
	if err := db.Create(&personalSession).Error; err != nil {
		t.Fatalf("create personal session: %v", err)
	}
	if err := db.Create(&secondSession).Error; err != nil {
		t.Fatalf("create second session: %v", err)
	}

	return workspaceBusinessFixture{
		User:              user,
		PersonalWorkspace: personalWorkspace,
		SecondWorkspace:   secondWorkspace,
		Ingredient:        ingredient,
		PersonalRecipe:    personalRecipe,
		SecondRecipe:      secondRecipe,
		PersonalPackage:   personalPackage,
		SecondPackage:     secondPackage,
		PersonalProduct:   personalProduct,
		SecondProduct:     secondProduct,
		PersonalClient:    personalClient,
		SecondClient:      secondClient,
		PersonalOrder:     personalOrder,
		SecondOrder:       secondOrder,
		PersonalSession:   personalSession,
		SecondSession:     secondSession,
	}
}

func runWorkspaceJSONRequest(
	userID uint,
	workspaceID uint,
	handler gin.HandlerFunc,
	method string,
	routePattern string,
	requestTarget string,
	body any,
) *httptest.ResponseRecorder {
	router := gin.New()
	router.Handle(method, routePattern, func(c *gin.Context) {
		c.Set("userID", userID)
		c.Set("workspaceID", workspaceID)
		handler(c)
	})

	var payload bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&payload).Encode(body); err != nil {
			panic(err)
		}
	}
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(method, requestTarget, &payload)
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(recorder, request)
	return recorder
}

func TestBusinessEndpointsFilterByWorkspace(t *testing.T) {
	fixture := setupWorkspaceBusinessTest(t)

	assertWorkspaceList(t, fixture, GetRecipes, "/recipes", fixture.PersonalWorkspace.ID, "recipes", fixture.PersonalRecipe.ID)
	assertWorkspaceList(t, fixture, GetPackages, "/packages", fixture.PersonalWorkspace.ID, "packages", fixture.PersonalPackage.ID)
	assertWorkspaceList(t, fixture, GetClients, "/clients", fixture.PersonalWorkspace.ID, "clients", fixture.PersonalClient.ID)
	assertWorkspaceList(t, fixture, GetProducts, "/products", fixture.PersonalWorkspace.ID, "products", fixture.PersonalProduct.ID)
	assertWorkspaceList(t, fixture, GetOrders, "/orders", fixture.PersonalWorkspace.ID, "orders", fixture.PersonalOrder.ID)
	assertWorkspaceList(t, fixture, GetCookingSessions, "/cooking_sessions", fixture.PersonalWorkspace.ID, "sessions", fixture.PersonalSession.ID)

	assertWorkspaceList(t, fixture, GetRecipes, "/recipes", fixture.SecondWorkspace.ID, "recipes", fixture.SecondRecipe.ID)
	assertWorkspaceList(t, fixture, GetPackages, "/packages", fixture.SecondWorkspace.ID, "packages", fixture.SecondPackage.ID)
	assertWorkspaceList(t, fixture, GetClients, "/clients", fixture.SecondWorkspace.ID, "clients", fixture.SecondClient.ID)
	assertWorkspaceList(t, fixture, GetProducts, "/products", fixture.SecondWorkspace.ID, "products", fixture.SecondProduct.ID)
	assertWorkspaceList(t, fixture, GetOrders, "/orders", fixture.SecondWorkspace.ID, "orders", fixture.SecondOrder.ID)
	assertWorkspaceList(t, fixture, GetCookingSessions, "/cooking_sessions", fixture.SecondWorkspace.ID, "sessions", fixture.SecondSession.ID)
}

func TestCrossWorkspaceReferencesAreRejected(t *testing.T) {
	fixture := setupWorkspaceBusinessTest(t)

	productWithOtherPackage := map[string]any{
		"name":       "Invalid package product",
		"price":      10,
		"cost":       4,
		"package_id": fixture.SecondPackage.ID,
		"recipe_ids": []uint{fixture.PersonalRecipe.ID},
	}
	response := runWorkspaceJSONRequest(fixture.User.ID, fixture.PersonalWorkspace.ID, CreateProduct, http.MethodPost, "/products", "/products", productWithOtherPackage)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("product with other workspace package status = %d body = %s", response.Code, response.Body.String())
	}

	productWithOtherRecipe := map[string]any{
		"name":       "Invalid recipe product",
		"price":      10,
		"cost":       4,
		"package_id": fixture.PersonalPackage.ID,
		"recipe_ids": []uint{fixture.SecondRecipe.ID},
	}
	response = runWorkspaceJSONRequest(fixture.User.ID, fixture.PersonalWorkspace.ID, CreateProduct, http.MethodPost, "/products", "/products", productWithOtherRecipe)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("product with other workspace recipe status = %d body = %s", response.Code, response.Body.String())
	}

	orderWithOtherClient := orderPayload(fixture.SecondClient.ID, fixture.PersonalProduct.ID)
	response = runWorkspaceJSONRequest(fixture.User.ID, fixture.PersonalWorkspace.ID, AddOrder, http.MethodPost, "/orders", "/orders", orderWithOtherClient)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("order with other workspace client status = %d body = %s", response.Code, response.Body.String())
	}

	orderWithOtherProduct := orderPayload(fixture.PersonalClient.ID, fixture.SecondProduct.ID)
	response = runWorkspaceJSONRequest(fixture.User.ID, fixture.PersonalWorkspace.ID, AddOrder, http.MethodPost, "/orders", "/orders", orderWithOtherProduct)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("order with other workspace product status = %d body = %s", response.Code, response.Body.String())
	}

	sessionPayload := map[string]any{
		"recipe_id": fixture.SecondRecipe.ID,
		"date":      time.Now(),
		"yield":     "1 kg",
	}
	response = runWorkspaceJSONRequest(fixture.User.ID, fixture.PersonalWorkspace.ID, CreateCookingSession, http.MethodPost, "/cooking_sessions", "/cooking_sessions", sessionPayload)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("session with other workspace recipe status = %d body = %s", response.Code, response.Body.String())
	}

	ingredientPayload := map[string]any{
		"ingredient_id": fixture.Ingredient.ID,
		"quantity":      "100",
		"unit":          "g",
	}
	response = runWorkspaceJSONRequest(
		fixture.User.ID,
		fixture.PersonalWorkspace.ID,
		AddIngredientToRecipe,
		http.MethodPost,
		"/recipes/:id/ingredients",
		"/recipes/"+uintToString(fixture.SecondRecipe.ID)+"/ingredients",
		ingredientPayload,
	)
	if response.Code != http.StatusNotFound {
		t.Fatalf("ingredient add to other workspace recipe status = %d body = %s", response.Code, response.Body.String())
	}
}

func TestDashboardUsesActiveWorkspace(t *testing.T) {
	fixture := setupWorkspaceBusinessTest(t)

	response := runWorkspaceRequest(fixture.User.ID, fixture.PersonalWorkspace.ID, GetDashboardData, http.MethodGet, "/dashboard", "/dashboard")
	if response.Code != http.StatusOK {
		t.Fatalf("dashboard status = %d body = %s", response.Code, response.Body.String())
	}
	var dashboard DashboardData
	if err := json.Unmarshal(response.Body.Bytes(), &dashboard); err != nil {
		t.Fatalf("decode dashboard: %v", err)
	}
	if dashboard.TotalRecipes != 1 || dashboard.TotalProducts != 1 || dashboard.TotalOrders != 1 {
		t.Fatalf("dashboard totals = recipes:%d products:%d orders:%d, want 1/1/1", dashboard.TotalRecipes, dashboard.TotalProducts, dashboard.TotalOrders)
	}
	if dashboard.PendingOrders != 0 {
		t.Fatalf("pending orders = %d, want 0", dashboard.PendingOrders)
	}
	if len(dashboard.RecentOrders) != 1 || dashboard.RecentOrders[0].ID != fixture.PersonalOrder.ID || dashboard.RecentOrders[0].TotalAmount != 24 {
		t.Fatalf("recent orders = %+v, want personal order total 24", dashboard.RecentOrders)
	}

	response = runWorkspaceRequest(fixture.User.ID, fixture.PersonalWorkspace.ID, GetProfitData, http.MethodGet, "/dashboard/profit", "/dashboard/profit")
	if response.Code != http.StatusOK {
		t.Fatalf("profit status = %d body = %s", response.Code, response.Body.String())
	}
	var profit ProfitData
	if err := json.Unmarshal(response.Body.Bytes(), &profit); err != nil {
		t.Fatalf("decode profit: %v", err)
	}
	if profit.TotalRevenue != 24 || profit.TotalCosts != 10 || profit.TotalProfit != 14 || profit.OrderCount != 1 {
		t.Fatalf("profit = %+v, want revenue 24 costs 10 profit 14 count 1", profit)
	}
}

func assertWorkspaceList(
	t *testing.T,
	fixture workspaceBusinessFixture,
	handler gin.HandlerFunc,
	path string,
	workspaceID uint,
	name string,
	wantID uint,
) {
	t.Helper()

	response := runWorkspaceRequest(fixture.User.ID, workspaceID, handler, http.MethodGet, path, path)
	if response.Code != http.StatusOK {
		t.Fatalf("%s status = %d body = %s", name, response.Code, response.Body.String())
	}
	var rows []struct {
		ID uint `json:"id"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &rows); err != nil {
		t.Fatalf("decode %s: %v", name, err)
	}
	if len(rows) != 1 || rows[0].ID != wantID {
		t.Fatalf("%s rows = %+v, want only id %d", name, rows, wantID)
	}
}

func orderPayload(clientID uint, productID uint) map[string]any {
	return map[string]any{
		"client_id": clientID,
		"status":    constants.OrderStatusNew,
		"items": []map[string]any{
			{
				"product_id": productID,
				"quantity":   1,
				"price":      10,
				"cost_price": 4,
			},
		},
	}
}
