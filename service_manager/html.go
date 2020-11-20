package service_manager

const html = `
<!doctype html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0"><meta http-equiv="X-UA-Compatible" content="ie=edge">
	<title>Snow Web Admin</title>
	<link rel="stylesheet" type="text/css" href="https://www.layuicdn.com/layui/css/layui.css" />
	<script src="https://cdn.bootcdn.net/ajax/libs/jquery/1.10.0/jquery.min.js"></script>
	<script src="https://cdn.bootcdn.net/ajax/libs/jsrender/0.9.71/jsrender.min.js"></script>
<!--	<script src="https://cdn.bootcdn.net/ajax/libs/layui/2.4.3/layui.all.min.js"></script>-->
</head>
<body>
	<div id="node_list">正在载入节点..</div>
	
</body>
<script type="text/x-jsrender" id="t-node_list">
<table class="layui-table">
<caption>节点列表</caption>
<tr>
	<th lay-data="{sort:true}">id</th>
	<th lay-data="{sort:true}">节点</td>
	<th lay-data="{sort:true}">宿主机</td>
	<th lay-data="{sort:true}">Go运行时</th>
	<th lay-data="{sort:true}">可启动服务</th>
	<th lay-data="{sort:true}">已启动服务</th>
</tr>
{{for Data}}
	<tr>
		<td>{{:#index + 1}}</td>
		<td>
			<table>
			<tr><td>节点名称</td> <td>{{:Name}}</td></tr>
			<tr><td>Snow版本</td> <td>{{:SnowVersion}}</td></tr>
			
			</table>
		</td>
		<td>
			<table>
			<tr><td>名称</td> <td>{{:Os.Hostname}}</td></tr>
			<tr><td>iP</td> <td>{{:Os.Ip}}</td></tr>
			</table>
		</td>
		<td>
		<table>
			<tr><td>编译环境</td> <td>{{:Runtime.GOOS}}/{{:Runtime.GOARCH}}</td></tr>
			<tr><td>Go版本</td> <td>{{:Runtime.GoVersion}}</td></tr>
			<tr><td>可使用CPU数</td> <td>{{:Runtime.NumCpu}}</td></tr>
			<tr><td>Goroutine数</td> <td>{{:Runtime.NumGoroutine}}</td></tr>
		</table>
		</td>
		<td>
			<table>
				<tr>
					<th>id</th> <th>名称</th> <th>操作</th>
				</tr>

				{{for Services}}
				<tr>
					<td>{{:#index + 1}}</td>
					<td><a href="" title="{{:Remark}}">{{:Name}}</a></td>
					<td> <button class="layui-btn" title="点击启动该服务" href="">启动</button> </td>
				</tr>
				{{/for}}
			</table>
		</td>
		<td>
			<table>
				<tr>
					<th>id</th> <th>名称</th> <th>启动时间</th> <th>操作</th>
				</tr>

				{{for Active}}
				<tr>
					<td>{{:#index + 1}}</td>
					<td><a href="">{{:Name}}</a></td>
					<td>{{:~formatUnix(MountTime)}}</td>
					<td> <button class="layui-btn" title="点击关闭该服务" href="">关闭</button> </td>
				</tr>
				{{/for}}
			</table>
		</td>
	</tr>
{{/for}}
</table>
</script>
<script>

$(document).ready(function() {
  $.views.helpers({
  	formatUnix: function (unix){
		var date = new Date(unix*1000);
		var Y = date.getFullYear()+1 + '-';
		var M = (date.getMonth()+1 < 10 ? '0'+(date.getMonth()+1) : date.getMonth()+1) + '-';
		var D = (date.getDate()+1  < 10 ? '0'+date.getDate()  : date.getDate())  + ' ';
		var h = date.getHours() + ':';
		var m = date.getMinutes() + ':';
		var s=  (date.getSeconds() < 10 ? '0'+ date.getSeconds()  : date.getSeconds());
		return Y+M+D+h+m+s;
	}
  })
  $.ajax({url: "/api/nodes", method:"post", success: function(res) {
      res = eval("("+res+")")
      console.log(res)
      // 输出模板
	  $("#node_list").html($.templates('#t-node_list')(res))
  }, error: function(res) {
    window.alert("列表加载失败！请检查本地网络");
  }})
})
</script>
</html>
`
