package storage

import (
	"fmt"

	"github.com/ProgrammistNik/WB-L3/l3.7_WarehouseControl/internal/model"
	"github.com/wb-go/wbf/dbpg"
)

type Storage struct {
	db *dbpg.DB
}

func New(db *dbpg.DB) *Storage {
	return &Storage{db: db}
}

func (st *Storage) ListItems() ([]model.Item, error) {
	rows, err := st.db.Master.Query(`SELECT id, name, quantity, updated_at FROM items ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []model.Item
	for rows.Next() {
		var it model.Item
		if err := rows.Scan(&it.ID, &it.Name, &it.Quantity, &it.UpdatedAt); err != nil {
			return nil, err
		}
		res = append(res, it)
	}
	return res, nil
}

func (st *Storage) CreateItem(name string, qty int, username string) error {
	// Устанавливаем правильную переменную для триггера
	if _, err := st.db.Master.Exec(fmt.Sprintf("SET \"myapp.current_user\" = '%s'", username)); err != nil {
		return err
	}
	_, err := st.db.Master.Exec(`INSERT INTO items(name, quantity) VALUES($1,$2)`, name, qty)
	return err
}

func (st *Storage) UpdateItem(id int, name string, qty int, username string) error {
	if _, err := st.db.Master.Exec(fmt.Sprintf("SET \"myapp.current_user\" = '%s'", username)); err != nil {
		return err
	}
	_, err := st.db.Master.Exec(`UPDATE items SET name=$1, quantity=$2 WHERE id=$3`,
		name, qty, id)
	return err
}

func (st *Storage) DeleteItem(id int, username string) error {
	if _, err := st.db.Master.Exec(fmt.Sprintf("SET \"myapp.current_user\" = '%s'", username)); err != nil {
		return err
	}
	_, err := st.db.Master.Exec(`DELETE FROM items WHERE id=$1`, id)
	return err
}

func (st *Storage) GetHistory(id int) ([]model.ItemHistory, error) {
	rows, err := st.db.Master.Query(`
        SELECT id, item_id, action, changed_by, old_data, new_data, changed_at 
        FROM item_history WHERE item_id=$1 ORDER BY changed_at DESC`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []model.ItemHistory
	for rows.Next() {
		var h model.ItemHistory
		if err := rows.Scan(
			&h.ID, &h.ItemID, &h.Action, &h.ChangedBy,
			&h.OldData, &h.NewData, &h.ChangedAt); err != nil {
			return nil, err
		}
		res = append(res, h)
	}
	return res, nil
}
