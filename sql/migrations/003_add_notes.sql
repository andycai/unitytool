-- 笔记分类表
CREATE TABLE note_categories (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL COMMENT '分类名称',
    description TEXT COMMENT '分类描述',
    parent_id BIGINT DEFAULT 0 COMMENT '父分类ID，0表示顶级分类',
    is_public TINYINT(1) DEFAULT 0 COMMENT '是否公开，0-否，1-是',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by BIGINT NOT NULL COMMENT '创建者ID',
    updated_by BIGINT NOT NULL COMMENT '更新者ID',
    INDEX idx_parent (parent_id),
    INDEX idx_public (is_public)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='笔记分类表';

-- 笔记表
CREATE TABLE notes (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(200) NOT NULL COMMENT '笔记标题',
    content TEXT NOT NULL COMMENT '笔记内容(Markdown)',
    category_id BIGINT NOT NULL COMMENT '分类ID',
    parent_id BIGINT DEFAULT 0 COMMENT '父笔记ID，0表示顶级笔记',
    is_public TINYINT(1) DEFAULT 0 COMMENT '是否公开，0-否，1-是',
    view_count INT DEFAULT 0 COMMENT '��览次数',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by BIGINT NOT NULL COMMENT '创建者ID',
    updated_by BIGINT NOT NULL COMMENT '更新者ID',
    INDEX idx_category (category_id),
    INDEX idx_parent (parent_id),
    INDEX idx_public (is_public),
    FOREIGN KEY (category_id) REFERENCES note_categories(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='笔记表';

-- 笔记访问权限表
CREATE TABLE note_permissions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    note_id BIGINT NOT NULL COMMENT '笔记ID',
    role_id BIGINT NOT NULL COMMENT '角色ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT NOT NULL COMMENT '创建者ID',
    UNIQUE KEY uk_note_role (note_id, role_id),
    FOREIGN KEY (note_id) REFERENCES notes(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='笔记访问权限表';

-- 分类访问权限表
CREATE TABLE category_permissions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    category_id BIGINT NOT NULL COMMENT '分类ID',
    role_id BIGINT NOT NULL COMMENT '角色ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT NOT NULL COMMENT '创建者ID',
    UNIQUE KEY uk_category_role (category_id, role_id),
    FOREIGN KEY (category_id) REFERENCES note_categories(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='分类访问权限表'; 