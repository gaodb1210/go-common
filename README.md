# go-common
积累一些常用的go公共库


## 1、command

- runner用来执行一些命令，可以设置超时时间。
- 可以使用SyncRunSample执行命令，忽视命令的输出。
- 也可以使用SynRun，传入stdoutWriter和stderrWriter，用来接收命令输出信息。
