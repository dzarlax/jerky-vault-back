# Migration Notes

## Order Status Update

### New Order Statuses (v2.0)

The order status system has been updated with new standardized statuses:

#### Old Statuses → New Statuses Mapping:
- `pending` → `new`
- `completed` → `finished`
- `Finished` → `finished`

#### New Status Flow:
1. **`new`** - Новый заказ (устанавливается по умолчанию)
2. **`in_progress`** - Заказ в процессе выполнения
3. **`ready`** - Заказ готов к выдаче/отправке
4. **`finished`** - Заказ завершен
5. **`canceled`** - Заказ отменен

#### Database Migration Required:
```sql
-- Update existing order statuses
UPDATE orders SET status = 'new' WHERE status = 'pending';
UPDATE orders SET status = 'finished' WHERE status = 'completed' OR status = 'Finished';
```

#### Code Changes Made:
- Added `constants/order_status.go` with status constants
- Updated `controllers/orders.go` to use new default status
- Updated `controllers/dashboard.go` to filter finished and canceled orders
- Updated `API_DOCUMENTATION.md` with new status examples

#### Dashboard Changes:
- Pending orders now exclude both `finished` and `canceled` statuses
- Order distribution shows all 5 new statuses

### Breaking Changes:
- API responses now return new status values
- Default order status changed from `"pending"` to `"new"`
- Dashboard filtering logic updated 