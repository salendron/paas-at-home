package main

type Permission struct {
	Key  string                 `json:"key"`
	Meta map[string]interface{} `json:"meta"`
}

type PermissionConfig struct {
	RootPermission *Permissions `json:"config"`
}

type Permissions struct {
	Key              string         `json:"key"`
	Name             string         `json:"name"`
	ChildPermissions []*Permissions `json:"child-permissions"`
}

func GetPermissionConfig() *PermissionConfig {
	p := &PermissionConfig{}

	// SUPER USER
	p.RootPermission = &Permissions{
		Key:              "ROOT",
		Name:             "Super User",
		ChildPermissions: make([]*Permissions, 0),
	}

	// DATA LOGGER
	dataLoggerPermissions := &Permissions{
		Key:  "DATA-LOGGER-ADMIN",
		Name: "Data Logger",
		ChildPermissions: []*Permissions{
			&Permissions{
				Key:              "DATA-LOGGER-READ",
				Name:             "Read",
				ChildPermissions: make([]*Permissions, 0),
			},
			&Permissions{
				Key:              "DATA-LOGGER-WRITE",
				Name:             "Write",
				ChildPermissions: make([]*Permissions, 0),
			},
		},
	}
	p.RootPermission.ChildPermissions = append(p.RootPermission.ChildPermissions, dataLoggerPermissions)

	// PUBSUB
	pubsubPermissions := &Permissions{
		Key:  "PUBSUB-ADMIN",
		Name: "PubSub",
		ChildPermissions: []*Permissions{
			&Permissions{
				Key:              "PUBSUB-MANAGE-TOPICS",
				Name:             "Manage Topics",
				ChildPermissions: make([]*Permissions, 0),
			},
			&Permissions{
				Key:              "PUBSUB-PUBLISH",
				Name:             "Publish",
				ChildPermissions: make([]*Permissions, 0),
			},
		},
	}
	p.RootPermission.ChildPermissions = append(p.RootPermission.ChildPermissions, pubsubPermissions)

	// IN MEMORY DB
	inMemoryDBPermissions := &Permissions{
		Key:  "IN-MEMORY-DB-ADMIN",
		Name: "In-Memory DB",
		ChildPermissions: []*Permissions{
			&Permissions{
				Key:              "IN-MEMORY-DB-READ",
				Name:             "Read",
				ChildPermissions: make([]*Permissions, 0),
			},
			&Permissions{
				Key:              "IN-MEMORY-DB-WRITE",
				Name:             "Write",
				ChildPermissions: make([]*Permissions, 0),
			},
		},
	}
	p.RootPermission.ChildPermissions = append(p.RootPermission.ChildPermissions, inMemoryDBPermissions)

	return p
}
