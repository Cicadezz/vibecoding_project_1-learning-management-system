const fs = require('fs');

const path = '.worktrees/codex-learning-growth-mvp/backend/internal/database/models/models_test.go';
let text = fs.readFileSync(path, 'utf8');
const oldImport = ['import (','\t\ testing\','','\t\github.com/glebarez/sqlite\','\t\gorm.io/gorm\',')',''].join('\n');
const newImport = ['import (','\t\ strings\','\t\testing\','','\t\gorm.io/driver/sqlite\','\t\gorm.io/gorm\',')',''].join('\n');
text = text.replace(oldImport, newImport);
const oldErr = ['\tif err != nil {','\t\t\tt.Fatalf(\" open sqlite: "%v\, err)','\t}',''].join('\n');
