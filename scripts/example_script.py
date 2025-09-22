#!/usr/bin/env python3
"""
示例Python脚本，用于测试RPA Middleware的Python脚本执行功能
"""

import os
import sys
import json
import time
from datetime import datetime

def main():
    """主函数"""
    print("Python脚本开始执行...")
    
    # 获取命令行参数
    args = sys.argv[1:] if len(sys.argv) > 1 else []
    print(f"接收到的参数: {args}")
    
    # 获取环境变量
    env_vars = dict(os.environ)
    print(f"环境变量数量: {len(env_vars)}")
    
    # 模拟一些处理逻辑
    print("正在处理数据...")
    time.sleep(2)  # 模拟处理时间
    
    # 创建text目录（如果不存在）
    text_dir = "text"
    if not os.path.exists(text_dir):
        os.makedirs(text_dir)
        print(f"创建目录: {text_dir}")
    
    # 生成markdown文件
    timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
    filename = f"output_{timestamp}.md"
    filepath = os.path.join(text_dir, filename)
    
    # 生成markdown内容
    content = f"""# Python脚本执行结果

## 执行时间
- 开始时间: {datetime.now().strftime("%Y-%m-%d %H:%M:%S")}
- 脚本路径: {os.path.abspath(__file__)}

## 执行参数
- 参数数量: {len(args)}
- 参数列表: {args}

## 环境信息
- Python版本: {sys.version}
- 工作目录: {os.getcwd()}
- 操作系统: {os.name}

## 处理结果
- 状态: 成功
- 输出文件: {filename}
- 文件路径: {filepath}

## 生成的数据
```json
{{
    "script_name": "{os.path.basename(__file__)}",
    "execution_time": "{datetime.now().isoformat()}",
    "arguments": {args},
    "output_file": "{filename}",
    "status": "completed"
}}
```

---
*此文件由Python脚本自动生成*
"""
    
    # 写入文件
    try:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(content)
        print(f"成功生成文件: {filepath}")
    except Exception as e:
        print(f"写入文件失败: {e}")
        return 1
    
    print("Python脚本执行完成!")
    return 0

if __name__ == "__main__":
    exit_code = main()
    sys.exit(exit_code)

