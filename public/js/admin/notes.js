function noteManagement() {
    return {
        categories: [],
        notes: [],
        roles: [],
        selectedNote: null,
        showCategoryModal: false,
        showNoteModal: false,
        categoryForm: {
            id: null,
            name: '',
            description: '',
            parent_id: 0,
            is_public: 0,
            role_ids: []
        },
        noteForm: {
            id: null,
            title: '',
            content: '',
            category_id: null,
            parent_id: 0,
            is_public: 0,
            role_ids: []
        },

        init() {
            this.loadCategories();
            this.loadNotes();
            this.loadRoles();
        },

        // 加载分类列表
        async loadCategories() {
            try {
                const response = await fetch('/admin/notes/categories');
                if (!response.ok) throw new Error('加载分类失败');
                this.categories = await response.json();
                // 初始化展开状态
                this.categories.forEach(category => {
                    category.expanded = false;
                });
            } catch (error) {
                console.error('Failed to load categories:', error);
                Alpine.store('notification').show('加载分类失败', 'error');
            }
        },

        // 加载笔记列表
        async loadNotes() {
            try {
                const response = await fetch('/admin/notes/tree');
                if (!response.ok) throw new Error('加载笔记失败');
                this.notes = await response.json();
                // 按分类组织笔记
                this.organizeNotes();
            } catch (error) {
                console.error('Failed to load notes:', error);
                Alpine.store('notification').show('加载笔记失败', 'error');
            }
        },

        // 加载角色列表
        async loadRoles() {
            try {
                const response = await fetch('/admin/roles');
                if (!response.ok) throw new Error('加载角色失败');
                this.roles = await response.json();
            } catch (error) {
                console.error('Failed to load roles:', error);
                Alpine.store('notification').show('加载角色失败', 'error');
            }
        },

        // 按分类组织笔记
        organizeNotes() {
            this.categories.forEach(category => {
                category.notes = this.notes.filter(note => note.category_id === category.id);
            });
        },

        // 切换分类展开状态
        toggleCategory(category) {
            category.expanded = !category.expanded;
        },

        // 选择笔记
        async selectNote(note) {
            try {
                const response = await fetch(`/admin/notes/${note.id}`);
                if (!response.ok) throw new Error('加载笔记详情失败');
                this.selectedNote = await response.json();
            } catch (error) {
                console.error('Failed to load note details:', error);
                Alpine.store('notification').show('加载笔记详情失败', 'error');
            }
        },

        // 创建笔记
        createNote() {
            this.noteForm = {
                id: null,
                title: '',
                content: '',
                category_id: this.categories[0]?.id,
                parent_id: 0,
                is_public: 0,
                role_ids: []
            };
            this.showNoteModal = true;
        },

        // 编辑笔记
        editNote(note) {
            this.noteForm = {
                id: note.id,
                title: note.title,
                content: note.content,
                category_id: note.category_id,
                parent_id: note.parent_id,
                is_public: note.is_public,
                role_ids: note.roles?.map(role => role.id) || []
            };
            this.showNoteModal = true;
        },

        // 保存笔记
        async saveNote() {
            if (this.selectedNote) {
                try {
                    const response = await fetch(`/admin/notes/${this.selectedNote.id}`, {
                        method: 'PUT',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify({
                            title: this.selectedNote.title,
                            content: this.selectedNote.content,
                            category_id: this.selectedNote.category_id,
                            parent_id: this.selectedNote.parent_id,
                            is_public: this.selectedNote.is_public ? 1 : 0,
                            role_ids: this.selectedNote.roles?.map(role => role.id) || []
                        }),
                    });

                    if (!response.ok) throw new Error('保存笔记失败');
                    
                    Alpine.store('notification').show('保存成功', 'success');
                    await this.loadNotes();
                } catch (error) {
                    console.error('Failed to save note:', error);
                    Alpine.store('notification').show('保存笔记失败', 'error');
                }
            }
        },

        // 提交笔记表单
        async submitNote() {
            try {
                const url = this.noteForm.id ? `/admin/notes/${this.noteForm.id}` : '/admin/notes';
                const method = this.noteForm.id ? 'PUT' : 'POST';

                const response = await fetch(url, {
                    method,
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        title: this.noteForm.title,
                        content: this.noteForm.content,
                        category_id: this.noteForm.category_id,
                        parent_id: this.noteForm.parent_id,
                        is_public: this.noteForm.is_public ? 1 : 0,
                        role_ids: this.noteForm.role_ids
                    }),
                });

                if (!response.ok) throw new Error('保存笔记失败');

                this.showNoteModal = false;
                Alpine.store('notification').show('保存成功', 'success');
                await this.loadNotes();
            } catch (error) {
                console.error('Failed to submit note:', error);
                Alpine.store('notification').show('保存笔记失败', 'error');
            }
        },

        // 删除笔记
        async deleteNote(note) {
            if (!confirm('确定要删除这个笔记吗？')) return;

            try {
                const response = await fetch(`/admin/notes/${note.id}`, {
                    method: 'DELETE',
                });

                if (!response.ok) throw new Error('删除笔记失败');

                if (this.selectedNote?.id === note.id) {
                    this.selectedNote = null;
                }

                Alpine.store('notification').show('删除成功', 'success');
                await this.loadNotes();
            } catch (error) {
                console.error('Failed to delete note:', error);
                Alpine.store('notification').show('删除笔记失败', 'error');
            }
        },

        // 编辑分类
        editCategory(category) {
            this.categoryForm = {
                id: category.id,
                name: category.name,
                description: category.description,
                parent_id: category.parent_id,
                is_public: category.is_public,
                role_ids: category.roles?.map(role => role.id) || []
            };
            this.showCategoryModal = true;
        },

        // 保存分类
        async saveCategory() {
            try {
                const url = this.categoryForm.id ? `/admin/notes/categories/${this.categoryForm.id}` : '/admin/notes/categories';
                const method = this.categoryForm.id ? 'PUT' : 'POST';

                const response = await fetch(url, {
                    method,
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        name: this.categoryForm.name,
                        description: this.categoryForm.description,
                        parent_id: this.categoryForm.parent_id,
                        is_public: this.categoryForm.is_public ? 1 : 0,
                        role_ids: this.categoryForm.role_ids
                    }),
                });

                if (!response.ok) throw new Error('保存分类失败');

                this.showCategoryModal = false;
                Alpine.store('notification').show('保存成功', 'success');
                await this.loadCategories();
            } catch (error) {
                console.error('Failed to save category:', error);
                Alpine.store('notification').show('保存分类失败', 'error');
            }
        },

        // 删除分类
        async deleteCategory(category) {
            if (!confirm('确定要删除这个分类吗？')) return;

            try {
                const response = await fetch(`/admin/notes/categories/${category.id}`, {
                    method: 'DELETE',
                });

                if (!response.ok) {
                    const error = await response.json();
                    throw new Error(error.message || '删除分类失败');
                }

                Alpine.store('notification').show('删除成功', 'success');
                await this.loadCategories();
            } catch (error) {
                console.error('Failed to delete category:', error);
                Alpine.store('notification').show(error.message || '删除分类失败', 'error');
            }
        },

        // Markdown 转 HTML
        markdownToHtml(markdown) {
            return marked.parse(markdown || '');
        }
    };
} 