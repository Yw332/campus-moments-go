#!/bin/bash

# 批量替换所有服务文件中的 DB 为 getDB()
cd /workspace/internal/service

# 需要替换的文件
files=("post_service.go" "comment_service.go" "like_service.go" "friend_service.go" "message_service.go" "tag_service.go")

for file in "${files[@]}"; do
    echo "正在修复 $file..."
    # 替换 DB. 为 getDB().
    sed -i 's/DB\./getDB()\./g' "$file"
    # 替换 DB.Preload 为 getDB().Preload
    sed -i 's/DB\.Preload/getDB().Preload/g' "$file"
    # 替换 DB.Model 为 getDB().Model
    sed -i 's/DB\.Model/getDB().Model/g' "$file"
    # 替换 DB.First 为 getDB().First
    sed -i 's/DB\.First/getDB().First/g' "$file"
    # 替换 DB.FirstOrCreate 为 getDB().FirstOrCreate
    sed -i 's/DB\.FirstOrCreate/getDB().FirstOrCreate/g' "$file"
    # 替换 DB.Save 为 getDB().Save
    sed -i 's/DB\.Save/getDB().Save/g' "$file"
    # 替换 DB.Create 为 getDB().Create
    sed -i 's/DB\.Create/getDB().Create/g' "$file"
    # 替换 DB.Update 为 getDB().Update
    sed -i 's/DB\.Update/getDB().Update/g' "$file"
    # 替换 DB.Updates 为 getDB().Updates
    sed -i 's/DB\.Updates/getDB().Updates/g' "$file"
    # 替换 DB.Delete 为 getDB().Delete
    sed -i 's/DB\.Delete/getDB().Delete/g' "$file"
    # 替换 DB.Where 为 getDB().Where
    sed -i 's/DB\.Where/getDB().Where/g' "$file"
    # 替换 DB.Begin 为 getDB().Begin
    sed -i 's/DB\.Begin/getDB().Begin/g' "$file"
    # 替换 DB.Commit 为 getDB().Commit
    sed -i 's/DB\.Commit/getDB().Commit/g' "$file"
    # 替换 DB.Rollback 为 getDB().Rollback
    sed -i 's/DB\.Rollback/getDB().Rollback/g' "$file"
    # 替换 DB.Transaction 为 getDB().Transaction
    sed -i 's/DB\.Transaction/getDB().Transaction/g' "$file"
    # 替换 DB.Count 为 getDB().Count
    sed -i 's/DB\.Count/getDB().Count/g' "$file"
    # 替换 DB.Pluck 为 getDB().Pluck
    sed -i 's/DB\.Pluck/getDB().Pluck/g' "$file"
    # 替换 DB.Find 为 getDB().Find
    sed -i 's/DB\.Find/getDB().Find/g' "$file"
    # 替换 DB.Order 为 getDB().Order
    sed -i 's/DB\.Order/getDB().Order/g' "$file"
    # 替换 DB.Limit 为 getDB().Limit
    sed -i 's/DB\.Limit/getDB().Limit/g' "$file"
    # 替换 DB.Offset 为 getDB().Offset
    sed -i 's/DB\.Offset/getDB().Offset/g' "$file"
    
    echo "✅ $file 修复完成"
done

echo "🎉 所有服务文件修复完成！"