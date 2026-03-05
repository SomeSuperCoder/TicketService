# Overdue Ticket Logic

## Centralized Calculation

The overdue ticket calculation logic is now centralized in the database view `v_ticket_overdue_status` (created in migration `005_create_overdue_view.sql`).

## How It Works

The view calculates for each ticket:
- `status_start_date`: When the ticket entered open/init status (or was created)
- `lost_days`: Days past the 7-day SLA
- `is_overdue`: Boolean flag if ticket is overdue

## SLA Configuration

The current SLA is **7 days** from when a ticket enters open/init status.

To change the SLA period, modify only the view definition in `005_create_overdue_view.sql`:
```sql
INTERVAL '7 days'  -- Change this value
```

## Queries Using This View

All monitoring queries now use this view:
1. `GetOverdueTickets` - Lists overdue tickets
2. `CountOverdueTickets` - Counts overdue tickets  
3. `GetDepartmentEfficiency` - Department overdue statistics

## Benefits

- **Single source of truth**: Change SLA in one place
- **Consistency**: All queries use identical logic
- **Performance**: View can be materialized if needed
- **Maintainability**: Easier to understand and modify
