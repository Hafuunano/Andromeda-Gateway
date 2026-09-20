module github.com/Hafuunano/Andromeda-Gateway

go 1.25.7

require (
	github.com/Hafuunano/Lucy v0.0.0
	github.com/Hafuunano/Protocol-ConvertTool v0.0.0
	github.com/joho/godotenv v1.5.1
	github.com/tencent-connect/botgo v0.0.0-20241218082132-fe31c0dfe469
	github.com/wdvxdr1123/ZeroBot v1.8.2
)

replace (
	github.com/Hafuunano/Core-SkillAction => ../Core-SkillAction
	github.com/Hafuunano/Lucy => ../Lucy
	github.com/Hafuunano/Plugin-Collections => ../Plugin-Collections
	github.com/Hafuunano/Protocol-ConvertTool => ../Protocol-ConvertTool
)
