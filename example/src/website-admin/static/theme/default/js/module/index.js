$(function(){
	$("input[id=module_delete]").click(function(){
		var modulename = $(this).parent().parent().children().first().text();

		var h3 = $("#deleteModal .modal-header h3");
		h3.html("是否确认删除 " + modulename + " ?");

		$("#deleteModal").modal("show");
		$("#module_delete_sure").click(function(){
			$.ajax({
				url: "delete_module.html",
				type: "POST",
				dataType: "json",
				data: "ajax=&modulename=" + modulename,
				success:function(json) {
	                window.location = "/module/";
	            }
			});
		})
    });

    $('input[id=module_stop]').click(function() {
		var modulename = $(this).parent().parent().children().first().text();

		var status = $(this).val();
		var h3 = $("#stopModal .modal-header h3");
		h3.html("是否确认" + status + " " + modulename + " ?");

		var input = document.getElementById("module_stop_sure");
		input.value = "确认" + status;

		$("#stopModal").modal("show");
		$("#module_stop_sure").click(function(){
			$.ajax({
				url: "stop_module.html",
				type: "POST",
				dataType: "json",
				data: "ajax=&modulename=" + modulename,
				success: function(json) {
				//	var stopButton = document.getElementById("module_stop");
				//	stopButton.value = "启用";
					window.location = "/module/";
				}
			});
		});
    });

    $('input[id=create_module]').click(function() {
    	$("#createModal").modal("show");
       	$("#btn_create").bind("click", function(){
		    var modulename = $("input[name=modulename]");
		    var modulenameHelp = modulename.next(".help-inline");
		    var modulenameParent = modulename.parents(".control-group");

		    var modulenameVal = modulename.val();

		    modulename.focus(function() {
		        modulenameParent.removeClass("error");
		        modulenameHelp.hide().html("");
		    });

		    if (!validateModulename(modulenameVal) ) {
		        modulenameParent.addClass("error");
		        modulenameHelp.html("模块名必须由1到20个下字符组成：大小写字母、数字或者下划线").show();
		        return;
		    }

		    $.ajax({
		    	url: "create_module.html",
		        type : "POST",
		        dataType : "json",
		        data: "ajax=&modulename="+modulenameVal,
		        success:function(json) {
		            console.log(json);
		            if (json.status == 0) {
		                modulenameParent.addClass("error");
		                modulenameHelp.html(json.message).show();
		            } else if(json.status == 1) {
		                window.location = "/module/";
		            }
		        }
		    });
    	});
	});
});