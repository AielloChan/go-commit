# Go Commit

## 介绍

一个用 Golang 编写的命令行工具，支持流程化的配置，并且有强大的字符串模板功能。

## 特性

- 无任何依赖包：二进制分发，下载后直接运行
- 原生跨平台：基于 Golang 的跨平台特性，能在 Window、macOS、Linux，不同 CPU 架构使用
- 全流程化的配置：按步骤配置，可实现根据上一步的结果跳转到指定下一步骤。可适配多种使用场景
- 强大的数据模板：所有的提示文本、选择文本、输出文本，都可以使用模板语法，使用 shell 的输出

## 配置文件示例

```json
{
  "stages": [
    {
      "label": "请选择改动的类型: ",
      "name": "type",
      "type": "select",
      "config": {
        "size": 4,
        "options": [
          {
            "label": "Feat\tA new feature",
            "value": "feat"
          },
          {
            "label": "Fix\tA bug fix",
            "value": "fix"
          },
          {
            "label": "Perf\tA code change that improves performance",
            "value": "perf"
          },
          {
            "label": "Test\tAdding missing tests or correcting existing tests",
            "value": "test"
          }
        ],
        "next": "scope"
      }
    },
    {
      "label": "改动的范围",
      "name": "scope",
      "type": "select",
      "config": {
        "size": 4,
        "options": [
          {
            "label": "首页\t\t首页以及该页面的更改",
            "value": "首页"
          },
          {
            "label": "人才主页\t人才主页相关的更改",
            "value": "人才主页"
          },
          {
            "label": "下载页面\t文档更新",
            "value": "下载页面"
          },
          {
            "label": "自定义\t手动编写 scope",
            "next": "customScope"
          }
        ]
      },
      "next": "title"
    },
    {
      "label": "请输入自定义的 scope",
      "name": "customScope",
      "type": "string",
      "config": {
        "min": 0,
        "max": 50
      },
      "next": "title"
    },
    {
      "label": "请输入标题",
      "name": "title",
      "type": "string",
      "config": {
        "min": 1,
        "max": 70
      }
    },
    {
      "label": "请输入更详细的描述",
      "name": "body",
      "type": "multiline",
      "config": {
        "min": 0,
        "max": 120
      },
      "next": "breaking"
    },
    {
      "label": "是否为破坏性修改",
      "name": "breaking",
      "type": "confirm",
      "config": {
        "default": false
      }
    },
    {
      "name": "checkStash",
      "type": "command",
      "config": {
        "cmd": "#![ `git diff --cached --name-only | wc -l` != 0 ]",
        "success": "submit",
        "failed": "noFile"
      }
    },
    {
      "label": "#!echo 加入的文件数量为 `git diff --cached --name-only | wc -l` 是否继续提交？",
      "name": "noFile",
      "type": "confirm",
      "next": "checkStash"
    },
    {
      "name": "submit",
      "type": "command",
      "config": {
        "cmd": "#$!git commit -m '{{.type}}({{.scope}}{{.customScope}}): {{.title}}\n\n{{.body}}'"
      }
    }
  ]
}
```

### 支持的步骤类型：

- select: 单选
- multi-select: 多选
- string: 文本输入
- multiline: 多行文本
- confirm: 是或否选择
- command: 执行命令

### 模板工具

所有的 `label` 字段和 `format` `success` 都支持模板工具，具体使用如下

- 以 `#!` 开头的会被当作 shell 进行执行，并且将其返回结果作为 label 字段的值，如 `"label": "#!echo '你好'"` 则会成为 `"label": "你好"`
- 以 `#$` 开头的会被当作模板进行执行，模板内能够进行插值，值为之前步骤中所获得的数据，并且将其返回结果作为 label 字段的值，如 `"label": "#$Type 是 {{.type}}"` 则会成为 `"label": "Type 是 fix"`
- 以 `#!$` 开头的则是先被当作 shell 执行，然后再被当作模板执行
- 以 `#$!` 开头的则是先被当作模板执行，然后再被当作 shell 执行，得到最终的结果作为 label 的值

所有的选项和步骤后，都可以定义 `next` 字段，它表示下一步该跳到哪个步骤执行，如 `改动的范围` 步骤中的最后一个选项，如果选择`自定义`那一项，则会直接跳转到 `请输入自定义的 scope` 这一步骤，实现流程可完全自定义化

## 测试使用

- 安装 golang 环境（自行 Google😂）
- clone 本项目到本地
- 切换到项目根目录
- 执行 `go mod download` 安装所需依赖包
- 执行 `go run main.go` 即可体验

## ToDo

- [x] 整理代码结构
- [x] 完善错误处理
- [ ] 增加测试
- [ ] 发布二进制包到包管理器（homebrew、yum、apt、WinGet）
- [ ] pull request [survey](https://github.com/AlecAivazis/survey) 的源码，使其多行编辑更友好
- [ ] pull request [survey](https://github.com/AlecAivazis/survey) 的源码，使其多语言提示更友好