// Copyright (c) 2026 Tencent Inc.
// SPDX-License-Identifier: Apache-2.0
//

package sandboxspec

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"github.com/tencentcloud/CubeSandbox/CubeMaster/pkg/base/db/models"
	sandboxtypes "github.com/tencentcloud/CubeSandbox/CubeMaster/pkg/service/sandbox/types"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupSpecStore(t *testing.T) *gorm.DB {
	t.Helper()
	database, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "spec.db")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)
	require.NoError(t, database.AutoMigrate(&models.SandboxSpec{}))
	require.NoError(t, database.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_sandbox_spec_sandbox_id ON t_cube_sandbox_spec (sandbox_id)").Error)

	dbMu.Lock()
	orig := db
	db = database
	dbMu.Unlock()
	t.Cleanup(func() {
		dbMu.Lock()
		db = orig
		dbMu.Unlock()
	})
	return database
}

func putSpec(t *testing.T, id, label string) {
	t.Helper()
	require.NoError(t, Put(context.Background(), id, &sandboxtypes.CreateCubeSandboxReq{
		InstanceType: "cubebox",
		Labels:       map[string]string{"app": label},
	}, PutOptions{}))
}
