'use strict';
(function () {
    //工具类
    var Permissions = {
        //获取地址栏参数
        initBtn: function (values) {
            $.ajax({
                url: "/system/btnpermissions",
                data:{btnIds:values},
                type:'GET',
                async:true,
                dataType: "json",
                success: function (items) {
                   if (items.ret==200){
                   	//console.info(items)
					   $.each(items.data,function (index, val) {
						   if (val.IsShow){
						       $("[data-permission='"+val.ControlUrl+"']").show();
						   }else{
                               $("[data-permission='"+val.ControlUrl+"']").remove();
						   }
                       })
				   }else{
                   	alert(items.msg);
				   }
                }
            });
        }
    }
   window.Permissions = Permissions;
}());