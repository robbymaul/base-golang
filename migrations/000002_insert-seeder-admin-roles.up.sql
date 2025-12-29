INSERT INTO admin_roles
    (id, code, name, description, permissions)
    OVERRIDING SYSTEM VALUE
VALUES (1, 'super_admin', 'super admin', '', '{}'),
       (2, 'admin', 'admin', '', '{}');
