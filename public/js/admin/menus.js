// Menu management functionality
function menuManagement() {
    return {
        menuTree: [],
        parentMenus: [],
        showModal: false,
        editMode: false,
        currentMenu: null,
        form: {
            parent_id: 0,
            name: '',
            path: '',
            icon: '',
            permission: '',
            sort: 0,
            is_show: true
        },
        loading: false,

        init() {
            this.fetchMenus();
        },

        async fetchMenus() {
            try {
                // 获取菜单树
                const treeResponse = await fetch('/api/menus/tree');
                if (!treeResponse.ok) throw new Error('获取菜单树失败');
                const treeData = await treeResponse.json();
                console.log('Menu Tree:', treeData); // Debug log
                this.menuTree = treeData;

                // 获取父级菜单列表（仅一级菜单）
                const response = await fetch('/api/menus');
                if (!response.ok) throw new Error('获取菜单列表失败');
                const menus = await response.json();
                this.parentMenus = menus.filter(menu => menu.parent_id === 0);
                console.log('Parent Menus:', this.parentMenus); // Debug log
            } catch (error) {
                console.error('Menu fetch error:', error); // Debug log
                Alpine.store('notification').show(error.message, 'error');
            }
        },

        createMenu() {
            this.editMode = false;
            this.currentMenu = null;
            this.form = {
                parent_id: 0,
                name: '',
                path: '',
                icon: '',
                permission: '',
                sort: 0,
                is_show: true
            };
            this.showModal = true;
        },

        editMenu(menu) {
            console.log('Editing menu:', menu); // Debug log
            this.editMode = true;
            this.currentMenu = menu;
            this.form = {
                parent_id: menu.parent_id,
                name: menu.name,
                path: menu.path,
                icon: menu.icon,
                permission: menu.permission,
                sort: menu.sort,
                is_show: menu.is_show
            };
            this.showModal = true;
        },

        closeModal() {
            this.showModal = false;
            this.editMode = false;
            this.currentMenu = null;
            this.form = {
                parent_id: 0,
                name: '',
                path: '',
                icon: '',
                permission: '',
                sort: 0,
                is_show: true
            };
        },

        async submitForm() {
            if (this.loading) return;
            this.loading = true;

            try {
                const url = this.editMode ? `/api/menus/${this.currentMenu.id}` : '/api/menus';
                const method = this.editMode ? 'PUT' : 'POST';
                
                console.log('Submitting form:', this.form); // Debug log
                const response = await fetch(url, {
                    method,
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(this.form)
                });

                if (!response.ok) {
                    const error = await response.json();
                    throw new Error(error.error || '操作失败');
                }

                Alpine.store('notification').show(
                    this.editMode ? '菜单更新成功' : '菜单创建成功',
                    'success'
                );
                this.closeModal();
                this.fetchMenus();
            } catch (error) {
                console.error('Submit error:', error); // Debug log
                Alpine.store('notification').show(error.message, 'error');
            } finally {
                this.loading = false;
            }
        },

        async deleteMenu(id) {
            if (!confirm('确定要删除这个菜单吗？如果是父级菜单，其下的子菜单也会被删除。')) return;

            try {
                console.log('Deleting menu:', id); // Debug log
                const response = await fetch(`/api/menus/${id}`, {
                    method: 'DELETE'
                });

                if (!response.ok) {
                    const error = await response.json();
                    throw new Error(error.error || '删除失败');
                }

                Alpine.store('notification').show('菜单删除成功', 'success');
                this.fetchMenus();
            } catch (error) {
                console.error('Delete error:', error); // Debug log
                Alpine.store('notification').show(error.message, 'error');
            }
        }
    }
} 