+ function ($, win) {
    "use strict";

    var GlobalFunc = new Global();
    var Tpl_Comment_list = _.template($('#TplCommentlist').html().replace(/\&lt\;/ig, '\<'));
    var Tpl_Province_list = _.template($('#TplProvinceList').html().replace(/\&lt\;/ig, '\<'));
    var Tpl_City_list = _.template($('#TplCityList').html().replace(/\&lt\;/ig, '\<'));
    var Tpl_tag_list = _.template($('#TplTagList').html().replace(/\&lt\;/ig, '\<'));
    var Id = location.href.split('?id=')[1];
    var name = location.href.match(/.*\/(.+).html/)[1];
    var $updateCustomerForm = $('#updateCustomerForm');
    var $updateActiveForm = $('#updateActiveForm');
    var $ModalAddOpport = $('#ModalAddOpport');
    var $addOpportNow = $('#addOpportNow');
    var $addOpportBtn = $('#addOpportBtn');
    var $addOpportForm = $('#addOpportForm');

    GlobalFunc.Focus('#addOpportNow', '#addOpportForm');


    win.addCommentBtn = function () {
        var p = {
            customer_id: Id,
            comment_area: $('#comment_area').val()
        }
        if (!$.trim(p.comment_area)) {
            GlobalFunc.Alert.info("备注不能为空!");
            return
        }

        $.post('/add_customer_comment', p, function (data) {
            if (!!data && data.code == 0) {
                GlobalFunc.Alert.success('增加成功 \^o^/');
                CustomInfoCtl.getCommentsUse();
            } else {
                GlobalFunc.Alert.fail('系统出错 X﹏X');
            }
        })
    }


    win.SaveOnClick = function () {
        if (GlobalFunc.CheckForm($updateCustomerForm) || GlobalFunc.CheckArea('#Province', '#City')) return;
        var Tags = [];
        $('#_updateTags span').each(function () {
            Tags.push($(this).data('id'))
        })

        // 发送请求
        var customer = {
            customerId: Id,
            entName: $('#EntName').val(),
            rtx_num: $('#Buin').val(),
            contacts: $('#Contacts').val(),
            phone: $('#ContactsPhone').val(),
            mobile: $('#Mobile').val(),
            mail: $('#ContactsMail').val(),
            city: $('#City option:checked').val(),
            qq: $('#QQ').val(),
            remarks: $('#Remarks').val(),
            tags: Tags
        }

        customer = CustomInfoCtl.checkNewInfo(customer);
        if (Object.getOwnPropertyNames(customer).length == 2) {
            CustomInfoCtl.noChange()
            return
        }
        $.post("/alert_customer", customer, function (data) {
            if (!!data && data.code == 0) {
                GlobalFunc.Alert.success('更新成功 \^o^/');
                CustomInfoCtl.getCustomInfo(Id);
                CustomInfoCtl.getCommentsUse();

            } else {
                GlobalFunc.Alert.fail('更新失败');
            }
        });

        CustomInfoCtl.noChange()
    }

    /*---------------------------------------------------------------------------------------------
    ---------------------------------------------------------------------------------------------*/

    var obj = {
        getCustomInfo: function (id) {
            var p = {
                customerId: id
            }
            var _this = this;
            var query;
            if (name == "get_emp_customer_detail") {
                query = '/getEmpCustomInfo'
            } else if (name == "get_agent_customer_detail") {
                query = 'getAgentCustomInfo'
            } else if (name == "get_super_customer_detail") {
                query = 'getSuperCustomInfo'
            }

            $.get(query, p, function (data) {
                if (!!data && data.code == -1) {
                    location.href = "404.html"
                    return
                }
                var _data = data;
                _data.customer.tags = [];
                _.map(data.customer.Tags, function (v) {
                    _data.customer.tags.push(v.tag_id);
                });
                _this.customInfo = _data.customer;
                $('#Buin').val(data.customer.rxt_num);
                $('#EntName').val(data.customer.entName);
                $('#Contacts').val(data.customer.contacts);
                $('#ContactsMail').val(data.customer.email);
                $('#ContactsPhone').val(data.customer.phone);
                $('#Mobile').val(data.customer.mobile);
                $('#QQ').val(data.customer.qq);
                $('#Remarks').val(data.customer.remarks);
                $('#Follow_time').val(data.customer.last_follow_time);

                $('#Active').val(data.customer.active)
                $('#Total').val(data.customer.EntYesterdayInfo.total)
                $('#IOS').val(data.customer.EntYesterdayInfo.iOS)
                $('#Android').val(data.customer.EntYesterdayInfo.android)
                $('#Pc').val(data.customer.EntYesterdayInfo.pc)

                _this.getAllProvince($.Deferred())
                    .done(function () {
                        $('#Province option').each(function () {
                            if ($(this).val() == data.customer.Province.province_id) {
                                $(this)[0].selected = true;
                                _this.getCity(data.customer.Province.province_id, 0, $.Deferred())
                                    .done(function () {
                                        $('#City option').each(function () {
                                            if ($(this).val() == data.customer.City.city_id) {
                                                $(this)[0].selected = true;
                                            }
                                        })
                                    })
                            }
                        })
                    });

                var h = Tpl_tag_list(data.customer);
                $('#_updateTags').html(h);
                _.map(data.customer.Tags, function (m) {
                    var c = '.tag_' + m.tag_id;
                    $('#_tagAll').find(c).parent().addClass('on');
                })
            })
        },

        getAllProvince: function (d) {
            $.post('/get_all_province', function (data) {
                if (!!data && data.code == 0) {
                    var h = Tpl_Province_list(data);
                    $('#Province').html(h);
                    d.resolve();
                }
            })
            return d.promise();
        },

        getCity: function (id, type, derfer) {
            var p = {
                province_id: id
            }
            $.post('/get_city', p, function (data) {
                if (!!data && data.code == 0) {
                    var h = Tpl_City_list(data);
                    type == 0 ? $('#City').html(h) : $('#_city').html(h);
                    derfer && derfer.resolve();
                }
            })
            var rsp = derfer ? derfer.promise() : 0;
            return rsp;
        },

        addOpport: function (Opport) {
            var _this = this;
            $.post('/add_opport', Opport, function (data) {
                if (!!data && data.code == 0) {
                    $ModalAddOpport.modal('hide');
                    $addOpportForm[0].reset();
                    GlobalFunc.Alert.success('转化成功 \^o^/');
                    _this.getCommentsUse();
                } else {
                    GlobalFunc.Alert.fail('系统出错 X﹏X');
                }
            })
        },

        updateComment: function (id, v) {
            var der = $.Deferred()
            $.post('/alert_customer_comment', {
                commentId: id,
                comments: v
            }, function (data) {
                if (!!data && data.code == 0) {
                    der.resolve(data)
                }
            })
            return der.promise()
        }
    }

    var CustomInfo = function () {
        this.lookInfo = {
            page: 1,
            type: 0
        }
        this.customInfo = null;
        this.getCustomInfo(Id);

        this.getCommentsUse = function (type, page, attr) {
            var p = {
                customerId: Id,
                type: type || 0,
                page: page || 1
            }
            attr = attr || 0;
            $.post('/get_sale_comments', p, function (data) {
                _.map(data.Comments, function (v) {
                    var _json = JSON.parse(v.comments)
                    v.comments = _json
                    v.text = ""
                    if (!v.type) {
                        v.text = _json["备注内容"]
                        return;
                    }
                    for (var i in _json) {
                        if (i == "action")
                            continue;
                        else
                            v.text += "[" + i + "] = " + _json[i] + "， "
                    }
                    v.text = v.text.substr(0, v.text.length - 2)
                });
                var h = Tpl_Comment_list(data);
                if (p.type == CustomInfoCtl.lookInfo.type && attr != 0) {
                    $('#CommentListBox').append(h);
                } else {
                    $('#CommentListBox').html(h);
                }
                CustomInfoCtl.lookInfo = p
            })
        };
        this.getCommentsUse();

        this.checkNewInfo = function (c) {
            c.entName == this.customInfo.entName && delete c.entName;
            c.rtx_num == this.customInfo.rxt_num && delete c.rtx_num;
            c.contacts == this.customInfo.contacts && delete c.contacts;
            c.phone == this.customInfo.phone && delete c.phone;
            c.mobile == this.customInfo.mobile && delete c.mobile;
            c.mail == this.customInfo.email && delete c.mail;
            c.city == this.customInfo.City.city_id && delete c.city;
            c.qq == this.customInfo.qq && delete c.qq;
            c.remarks == this.customInfo.remarks && delete c.remarks;
            if (c.tags.toString() == this.customInfo.tags.toString()) {
                c.tagsAlert = 0;
                delete c.tags;
            } else {
                c.tagsAlert = 1;
            }
            return c
        };

        this.noChange = function () {
            $('#editBtn').attr("disabled", false);
            $('#saveBtn').attr("disabled", true);

            $('#EntName').attr("disabled", true);
            $('#Buin').attr("disabled", true);
            $('#Contacts').attr("disabled", true);
            $('#ContactsPhone').attr("disabled", true);
            $('#Mobile').attr("disabled", true);
            $('#ContactsMail').attr("disabled", true);
            $('#QQ').attr("disabled", true);
            $('#Province').attr("disabled", true);
            $('#City').attr("disabled", true);
            $('#Remarks').attr("disabled", true);

            $('#_tagAll').addClass('hide');
        }
    }
    CustomInfo.prototype = obj;
    var CustomInfoCtl = new CustomInfo();



    /*--------------------------------------------------------------------------------------------
    --------------------------------------------------------------------------------------------*/

    $('.modal').modal({
        backdrop: 'static',
        keyboard: false,
        show: false
    });
    $('#Province').change(function () {
        var id = $(this).find('option:checked').val();
        if (!id) return;
        CustomInfoCtl.getCity(+id, 0)
    });

    $addOpportNow.click(function () {
        $ModalAddOpport.modal('show');
        $('#OpportTitle').val($('#EntName').val());
    });

    $addOpportBtn.on('click', function () {
        if (GlobalFunc.CheckForm($addOpportForm)) {
            return
        };
        var Opport = new Object();
        Opport.title = $('#OpportTitle').val();
        Opport.content = $('#Content').val();
        Opport.customerId = Id;
        Opport.status = $('#Stage').find('option:checked').val();
        Opport.quota = $('#Quota').val();
        Opport.estimatetime = $('#estimate_time').val();
        Opport.real_time = $('#real_time').val();

        var tags = [];
        $('#tagAll').find('.tag.on').each(function () {
            tags.push($(this).find('span').data('id'));
        });
        Opport.tags = tags;

        CustomInfoCtl.addOpport(Opport);
    });

    $('body').on('click', '.recordBtn', function () {
        CustomInfoCtl.getCommentsUse(1);
    });

    $('body').on('click', '.updateCommentBtn', function () {
        var textarea = $(this).parents('.panel').find('textarea');
        textarea.attr("readonly") == "readonly" ? textarea.removeAttr("readonly") : textarea.attr("readonly", true);
    });

    $('body').on('click', '.saveCommentBtn', function () {
        var dom = $(this).parents('.panel');
        var id = +dom.find('.comment_id').text();
        var value = dom.find('textarea').val()
        CustomInfoCtl.updateComment(id, value).then(function (data) {
            dom.find('textarea').attr("readonly", true);
        });
    });

    $('body').on('click', '.lookMore', function () {
        CustomInfoCtl.getCommentsUse(CustomInfoCtl.lookInfo.type, CustomInfoCtl.lookInfo.page + 1, 1);
    })

}(jQuery, window);