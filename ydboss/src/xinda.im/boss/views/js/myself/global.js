var obj = {
    pageTpl: "<%if(record>0){%>" +
        "<p class='text-muted font-12'>共<span class='text-primary'><%=record%></span>条记录 , 每页<span class='text-primary'><%=page_size%></span>条&nbsp;&nbsp;<%if(page_total==0){%>0<%}else{%><%=page%><%}%>/<span class='total'><%=page_total%></span></p>" +
        "<%if(page_total>1){%>" +
        "<ul class='pagination'>" +
        "<li><a class='trun-left' href='#'>«</a></li>" +
        "<%if(page_total<6){%>" +
        "<%for(var i=0;i<page_total;i++){%>" +
        "<li class='<%if((i+1)==page){%>active<%}%>'><a class='branch-btn' href='#'><%=i+1%></a></li>" +
        "<%}%>" +
        "<%}else{%>" +
        "<%for(var i=0;i<5;i++){%>" +
        "<%if(page>(page_total-2)){%>" +
        "<li class='<%if(i==(4-(page_total-page))){%>active<%}%>'><a class='branch-btn' href='#'><%=page_total-(4-i)%></a></li>" +
        "<%}else if(page<3){%>" +
        "<li class='<%if((i+1)==page){%>active<%}%>'><a class='branch-btn' href='#'><%=5-(4-i)%></a></li>" +
        "<%}else{%>" +
        "<li class='<%if(i==2){%>active<%}%>'><a class='branch-btn' href='#'><%=page-2+i%></a></li>" +
        "<%}%>" +
        "<%}%>" +
        "<%}%>" +
        "<li><a class='turn-right' href='#'>»</a></li>" +
        "</ul>" +
        "<%}%>" +
        "<%}%>",

    tagTpl: "<button class='btn btn-success tagAddBtn' type='button'><i class='fa fa-plus'></i> 添加标签</button>" +
        "<div class='col-sm-12 tagBox tagsinput'>" +
        "<i class='closeTagBox'></i>" +
        "<%for(var i =0;i<list.length;i++){%><span class='tag'>" +
        "<span data-id='<%=list[i].tag_id%>' class='tag_<%=list[i].tag_id%>'><%=list[i].Name%></span>" +
        "</span><%}%>" +
        "</div>",

    updatetagTpl: "<%if(list.length==0){%>" +
        "<%}else{%>" +
        "<%for(var i=0;i<list.length;i++){%>" +
        "<span data-id='<%=list[i].tag_id%>' class='label label-info'><%=list[i].Name%></span><%}%><%}%>",

    Assign_status: {
        "0": "未跟进"
    },

    OpportStatus: {
        "0": "打开",
        "1": "关闭"
    },

    OpportStage: {
        "4": "商务失败（0%）",
        "0": "初步交流（10%）",
        "1": "需求沟通（30%）",
        "2": "商务沟通（50%）",
        "3": "签约交款（100%）"
    },

    Regexp: {
        "noNull": {
            reg: /\S/,
            tip: "输入值不能为空!"
        },

        "email": {
            reg: /^$|^(\w)+(\.\w+)*@(\w)+((\.\w+)+)$/,
            tip: "请输入正确的邮箱格式！"
        },

        "Phone": {
            reg: /^$|^1[0-9]{10}$/,
            tip: "请输入正确的手机格式！"
        },

        "RTXNumber": {
            reg: /^$|^[1-9][0-9]{7}$/,
            tip: "请输入正确的RTX总机号格式！"
        },
        
        "psw": {
            reg: /^\S{6,12}$/,
            tip: "请输入6-12位任意数密码！"
        }
    },
    
    Menu: {
        "get_emp_customer_detail": "sale_customer",
        "get_agent_customer_detail" : "agent_customer",
        "get_super_customer_detail" : "all_customer",
        "get_emp_biz_detail": "opportunities",
        "employee_info": "agent_employee",
        "get_agent_biz_detail" : "agent_biz",
        "sale_customer_search": "sale_customer",
        "agent_customer_search": "agent_customer",
        "all_customer_search": "all_customer",
        "sale_opport_search": "opportunities",
        "agent_opport_search": "agent_biz"
    },
    
    UpdatePsw : function (p) {
        var _this = this;
        $.post('/update_account_psw', p, function (data) {
            if (!!data && data.code == 0) {
                $('#ModalUpdatePsw').modal('hide');
                _this.Alert.success('修改成功 \^o^/');
            } else {
                _this.Alert.fail('系统出错 X﹏X');
            }
        })
    },

    MenuAct: function () {
        var hrefs = location.href.match(/.*\/(.+).html/);
        if (!hrefs) return;
        var h = hrefs[1];
        h in this.Menu && (h = this.Menu[h]);
        var href = h + '.html';
        $('.left-side').find('.menu-list').find('a').each(function () {
            var _this = $(this);
            if (_this.attr('href') == href) {
                _this.parent('li').addClass('active');
                _this.parents('.menu-list').addClass('nav-active');
            }
        })
    },

    CheckForm: function (form) {
        var that = this;
        var rsp = 0;
        form.find('input').each(function () {
            var _this = $(this);
            var value = _this.val();
            var parentDom = _this.parents('.form-group').first();
            var data = _this.data('reg');
            var regs = data ? data.split(',') : [];
            _.forEach(regs, function (v) {
                v = $.trim(v);
                if (!that.Regexp[v].reg.test(value)) {
                    parentDom.addClass('has-error');
                    _this.siblings('p.help-block').text(that.Regexp[v].tip);
                    rsp++;
                    return false
                } else {
                    parentDom.removeClass('has-error');
                    _this.siblings('p.help-block').text('');
                }
            })
        })
        return rsp
    },

    Alert: {
        success: function (text) {
            var myToast = $.toast({
                heading: '操作成功',
                text: text,
                icon: 'success', //error,info,warning
                showHideTransition: 'fade', //slide，plain 
                loaderBg: '#5aa95a',
                //allowToastClose: false,
                stack: 4, //最大提示款存在数
                //hideAfter: true,   //false不会自动消失 ，填毫秒控制（5000）
                //position: 'top-right',
                bgColor: '#5cb85c',
                //textColor: 'white'
                position: {
                    right: 15,
                    top: 65
                }
            })
        },

        fail: function (text) {
            var myToast = $.toast({
                heading: '操作失败',
                text: text,
                icon: 'error',
                showHideTransition: 'fade',
                loaderBg: '#9EC600',
                stack: 4,
                position: {
                    right: 15,
                    top: 65
                }
            })
        },

        info: function (text) {
            var myToast = $.toast({
                heading: '提示',
                text: text,
                icon: 'info',
                showHideTransition: 'fade',
                loaderBg: '#9EC600',
                stack: 4,
                position: {
                    right: 15,
                    top: 65
                }
            })
        },

        warn: function (text) {
            var myToast = $.toast({
                heading: '警告',
                text: text,
                icon: 'warning',
                showHideTransition: 'fade',
                loaderBg: '#9EC600',
                stack: 4,
                position: {
                    right: 15,
                    top: 65
                }
            })
        }
    },

    Focus: function (btn, form) {
        $('body').on('click', btn, function (e) {
            var time = setTimeout(function () {
                $(form).find('input').not('.no-focus').first()[0].focus();
            }, 500)
        })
    },

    CheckArea: function (pro, city) {
        if (!$(pro).find('option:checked').val()) {
            $(pro).parents('.form-group').first().addClass('has-error');
            this.Alert.warn('请选择省份跟城市');
            return true
        }
        if (!$(city).find('option:checked').val()) {
            $(city).parents('.form-group').first().addClass('has-error');
            this.Alert.warn('请选择省份跟城市');
            return true
        }
        $(pro).parents('.form-group').first().removeClass('has-error');
        return false
    },

    CheckChoose: function (ids) {
        var g = 0;
        var dom = 'input[name="' + ids + '"]:checked';
        $(dom).each(function () {
            g++;
            return
        });
        return !g
    },

    TagInit: function () {
        var dom = $('#tagAll');
        var _dom = $('#_tagAll');
        var _this = this;
        if (dom.length == 0 && _dom.length == 0) return;

        $.post('/get_all_tag', function (data) {
            if (!!data && data.code == 0) {
                var h = _.template(_this.tagTpl)(data);
                dom.html(h);
                _dom.length != 0 && _dom.html(h);
            }
        })
    },

    CheckEntExist: function (name, der) {
        $.post("/customer_exists",{name: name}, function (data) {
            if (!!data && data.code == 0) {
                der.resolve()
            } else {
                der.reject()
            }
        })
        return der.promise()
    },

    _checkTable: function () {
        $('.alert .checkList').each(function () {
            var _this = $(this);
            var key = _this.data('key');
            var dom = $('.table th[data-key=' + key + '],.table td[data-key=' + key + ']');
            if (_this.hasClass('on')) {
                dom.removeClass('hide');
            } else {
                dom.addClass('hide');
            }
        })
    },

    _checkAll: function (dom) {
        var input = dom.siblings('.alert-title').find('.checkAll');
        dom.find('.checkList.on').length > 0 ? input.addClass('on') : input.removeClass('on');
    },

    _showTagsList: function (type) {
        var tags = {
            list: []
        };
        var dom = type == 'add' ? $('#tagAll .tag.on') : $('#_tagAll .tag.on');
        var box = type == 'add' ? $('#_addTags') : $('#_updateTags');
        dom.each(function () {
            var tag = {};
            tag.tag_id = $(this).find('span').data('id');
            tag.Name = $(this).find('span').text();
            tags.list.push(tag)
        });
        var h = _.template(this.updatetagTpl)(tags);
        box.html(h);
    },

    BindEventInit: function () {
        var That = this;
        $('.modal').modal({
            backdrop: 'static',
            keyboard: false,
            show: false
        });
        
        $('body').on('click', '.table tr:gt(0) input,.table tr:gt(0) button', function (e) {
            e.stopPropagation();
        });
        
        $('body').on('click', '.table tr:gt(0)', function () {
            if ($(this).hasClass('noHref')) {
                return
            }
            
            location.href = $(this).children('td').eq(1).find('a').attr('href');
        });

        $('body').on('click', '.checkList', function () {
            var _this = $(this);
            var key = _this.data('key');
            var dom = $('.table th[data-key=' + key + '],.table td[data-key=' + key + ']');
            _this.toggleClass('on');
            That._checkTable();
            That._checkAll(_this.parent('.alert'));
        });
        
        That.Focus('#changePsw', '#updatePswForm');
        $('#changePsw').click(function () {
            $('#ModalUpdatePsw').modal('show');
        });
        
        $('#UpdatePswBtn').click(function () {
            if (That.CheckForm($('#updatePswForm'))) return;
            var p = {
                new_psw: $('#newPsw').val()
            }
            That.UpdatePsw(p);
        });

        $('body').on('click', '.alert-title', function () {
            var This = $(this);
            var _this = This.find('.checkAll');
            _this.toggleClass('on');
            if (_this.hasClass('on')) {
                This.siblings('.alert').find('.checkList').addClass('on');
            } else {
                This.siblings('.alert').find('.checkList').removeClass('on');
            }
            That._checkTable();
        });

        $('body').on('click', '.tagAddBtn', function () {
            $(this).siblings('.tagBox').toggle();
        });

        $('body').on('click', '.closeTagBox', function () {
            $(this).parent('.tagBox').hide();
        });

        $('body').on('click', '#tagAll .tag', function () {
            $(this).toggleClass('on');
            That._showTagsList('add');
        });

        $('body').on('click', '#_tagAll .tag', function () {
            $(this).toggleClass('on');
            That._showTagsList('update');
        });
        
        $('.close').click(function () {
            var form = $(this).parent('.modal-header').siblings('.modal-body').find('form');
            form[0].reset();
            form.find('.form-group').removeClass('has-error');
            form.find('.help-block').html('');
        })
    }
}


var Global = function () {
    this.MenuAct();
    this.BindEventInit();
    this.TagInit();
}
Global.prototype = obj;