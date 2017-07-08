"user strict";
(function ($, win) {
    $('#Username').focus();
    var GlobalFunc = new Global();
    
    win.Post = function () {
        var name = document.getElementById('Username').value;
        var password = document.getElementById('password').value;
        if (!$.trim(name)) {
             $.toast({
                    heading: 'Info',
                    icon: 'info',
                    text: '账号不能为空'
                })
            return false
        }
        if (!$.trim(password)) {
             $.toast({
                    heading: 'Info',
                    icon: 'info',
                    text: '密码不能为空'
                })
            return false
        }
        var user = {
            name: name,
            password: password
        }
        $.post('/login', user, function (data) {
            if (!!data && data.status == 0) {
                data.href == 0 && (location.href = 'agent.html');
                data.href == 1 && (location.href = 'agent_employee.html');
                data.href == 2 && (location.href = 'sale_customer.html');
            } else if (!!data && data.status == 1) {
                $.toast({
                    heading: 'Error',
                    icon: 'error',
                    text: '账号不存在'
                })
            } else if (!!data && data.status == 2) {
                $.toast({
                    heading: 'Error',
                    icon: 'error',
                    text: '密码错误'
                })
            } else {
                $.toast({
                    heading: 'Error',
                    icon: 'error',
                    text: '系统出错'
                })
            }
        })
        return false;
    }
})(jQuery, window);