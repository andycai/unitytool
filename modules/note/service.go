package note

import (
	"github.com/andycai/unitool/models"
	"gorm.io/gorm"
)

// GetNoteByID 获取笔记
func GetNoteByID(id uint) (*models.Note, error) {
	var note models.Note
	if err := app.DB.Preload("Category").Preload("Roles").First(&note, id).Error; err != nil {
		return nil, err
	}
	return &note, nil
}

// GetCategoryByID 获取分类
func GetCategoryByID(id uint) (*models.NoteCategory, error) {
	var category models.NoteCategory
	if err := app.DB.Preload("Roles").First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// ListNotes 获取笔记列表
func ListNotes() ([]models.Note, error) {
	var notes []models.Note
	if err := app.DB.Preload("Category").Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

// ListCategories 获取分类列表
func ListCategories() ([]models.NoteCategory, error) {
	var categories []models.NoteCategory
	if err := app.DB.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// ListPublicNotes 获取公开笔记列表
func ListPublicNotes() ([]models.Note, error) {
	var notes []models.Note
	if err := app.DB.Where("is_public = ?", 1).Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

// ListPublicCategories 获取公开分类列表
func ListPublicCategories() ([]models.NoteCategory, error) {
	var categories []models.NoteCategory
	if err := app.DB.Where("is_public = ?", 1).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// HasCategoryChildren 检查分类是否有子分类
func HasCategoryChildren(id uint) (bool, error) {
	var count int64
	if err := app.DB.Model(&models.NoteCategory{}).Where("parent_id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// HasCategoryNotes 检查分类是否有笔记
func HasCategoryNotes(id uint) (bool, error) {
	var count int64
	if err := app.DB.Model(&models.Note{}).Where("category_id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// IncrementNoteViewCount 增加笔记浏览次数
func IncrementNoteViewCount(note *models.Note) error {
	return app.DB.Model(note).UpdateColumn("view_count", note.ViewCount+1).Error
}

// CreateNote 创建笔记
func CreateNote(note *models.Note, roleIDs []uint) error {
	return app.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(note).Error; err != nil {
			return err
		}

		if len(roleIDs) > 0 {
			var roles []models.Role
			if err := tx.Find(&roles, roleIDs).Error; err != nil {
				return err
			}
			if err := tx.Model(note).Association("Roles").Replace(roles); err != nil {
				return err
			}
		}

		return nil
	})
}

// UpdateNote 更新笔记
func UpdateNote(note *models.Note, roleIDs []uint) error {
	return app.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(note).Error; err != nil {
			return err
		}

		if len(roleIDs) > 0 {
			var roles []models.Role
			if err := tx.Find(&roles, roleIDs).Error; err != nil {
				return err
			}
			if err := tx.Model(note).Association("Roles").Replace(roles); err != nil {
				return err
			}
		} else {
			if err := tx.Model(note).Association("Roles").Clear(); err != nil {
				return err
			}
		}

		return nil
	})
}

// DeleteNote 删除笔记
func DeleteNote(note *models.Note) error {
	return app.DB.Delete(note).Error
}

// CreateCategory 创建分类
func CreateCategory(category *models.NoteCategory, roleIDs []uint) error {
	return app.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(category).Error; err != nil {
			return err
		}

		if len(roleIDs) > 0 {
			var roles []models.Role
			if err := tx.Find(&roles, roleIDs).Error; err != nil {
				return err
			}
			if err := tx.Model(category).Association("Roles").Replace(roles); err != nil {
				return err
			}
		}

		return nil
	})
}

// UpdateCategory 更新分类
func UpdateCategory(category *models.NoteCategory, roleIDs []uint) error {
	return app.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(category).Error; err != nil {
			return err
		}

		if len(roleIDs) > 0 {
			var roles []models.Role
			if err := tx.Find(&roles, roleIDs).Error; err != nil {
				return err
			}
			if err := tx.Model(category).Association("Roles").Replace(roles); err != nil {
				return err
			}
		} else {
			if err := tx.Model(category).Association("Roles").Clear(); err != nil {
				return err
			}
		}

		return nil
	})
}

// DeleteCategory 删除分类
func DeleteCategory(category *models.NoteCategory) error {
	return app.DB.Delete(category).Error
}
