"user struct";
(function (Dom) {

    /*声明*/
    var GlobalFunc = new Global();
    var AccountCtl = new Account();
    var Tpl_page = _.template(GlobalFunc.pageTpl),
        Tpl_Account_list = _.template($('#TplAccountList').html().replace(/\&lt\;/ig, '\<')),
        $addAccountBtn = $('#addAccountBtn'),
        $searchAccountBtn = $('#searchAccountBtn'),
        $delAccountBtn = $('#DelAccountBtn'),
        $Account_page = $('#Account_page'),
        $AgentList = $('#AgentList'),
        $AllocationBtn = $('#AllocationBtn'),
        $addAccountForm = $('#addAccountForm'),
        $ModalAddAccount = $('#ModalAddAccount'),
        $ModalChangeAccount = $('#ModalChangeAccount'),
        $delAccountNow = $('#delAccountNow'),
        $ModalDelAccount = $('#ModalDelAccount'),
        $changeAccountBtn = $('#changeAccountBtn'),
        $changeAccountForm = $('#changeAccountForm'),
        $addAccountNow = $('#addAccountNow'),
        $ModalResetPsw = $('#ModalResetPsw'),
        $resetPswBtn = $('#resetPswBtn')

    var That;

    GlobalFunc.Focus('#addAccountNow', '#addAccountForm');
    GlobalFunc.Focus('.changeAccountBtn', '#changeAccountForm');
    GlobalFunc.Focus('.updatePsw', '#resetPswForm');


    /*---------------------------------------------------------------------------------------------------
    ---------------------------------------------------------------------------------*/

    function Account() {
        this.Account = {};
        this.page_info = null;
        this.keyOn = null;
        this.key = null;
        this.list_length = 0;
        this.filter_info = {};
        this.sort = null;
        That = this;

        this.showAccountList = function (data) {
            var h_table = Tpl_Account_list(data);
            var h_page = Tpl_page(data.page_info);
            $('#AccountListBox').html(h_table);
            $Account_page.html(h_page);
        }

        this.changeAccount = function () {
            this.Account.AccountId = Dom.find('#_UserName').val();
            this.Account.Name = Dom.find('#_Name').val();
            this.Account.Gender = Dom.find('#_Sex').find('option:checked').val();
            this.Account.Mail = Dom.find('#_ContactsMail').val();
            this.Account.Mobile = Dom.find('#_ContactsPhone').val();
        }
        this.Flip = function (page) {
            switch (this.keyOn) {
            case 1:
                this.searchAccount(this.key, page);
                break;
            case 2:
                var o = this.sort;
                o.page = page;
                this.sortByKey(o);
                break;
            default:
                this.getAccountList(page)
            }
        }
    }

    Account.prototype.getAccountList = function (page) {
        var p = {
            page: page
        }
        $.post('/employee.list', p, function (data) {
            if (!!data && data.code == 0) {
                That.page_info = data.page_info;
                That.keyOn = 0;
                That.showAccountList(data);
            }
        })
    }

    Account.prototype.addAccount = function (Account) {
        $.post('/add_employee', Account, function (data) {
            if (!!data && data.code == 0) {
                That.getAccountList(1);
                $addAccountForm[0].reset();
                $ModalAddAccount.modal('hide');
                GlobalFunc.Alert.success('新增成功 \^o^/');
            } else if (!!data && data.code == 3) {
                GlobalFunc.Alert.info('改账号已存在，请重新输入新账号 X﹏X');
            } else {
                GlobalFunc.Alert.fail('系统出错 X﹏X');
            }
        })
    }

    Account.prototype.updateAccount = function (Account) {
        $.post('/alert_employee', Account, function (data) {
            if (!!data && data.code == 0) {
                That.getAccountList(1);
                $ModalChangeAccount.modal('hide');
                GlobalFunc.Alert.success('修改成功 \^o^/');
            } else {
                GlobalFunc.Alert.fail('系统出错 X﹏X');
            }
        })
    }

    Account.prototype.deleteAccount = function (ids) {
        var p = {
            AccountIds: ids
        }
        $.post('/del_employee', p, function (data) {
            if (!!data && data.code == 0) {
                That.getAccountList(1);
                GlobalFunc.Alert.success('删除成功 \^o^/');
            } else {
                GlobalFunc.Alert.fail('系统出错 X﹏X');
            }
        })
    }

    Account.prototype.searchAccount = function (key, page) {
        var p = {
            keyWord: key,
            page: page || 1
        }
        $.post('/keyWord_employee', p, function (data) {
            if (!!data && data.code == 0) {
                That.page_info = data.page_info;
                That.key = key;
                That.keyOn = 1;
                That.showAccountList(data);
            }
        })
    }

//    Account.prototype.accountFilter = function (o) {
//        $.post('/account_filter', o, function (data) {
//            if (!!data && data.code == 0) {
//                That.page_info = data.page_info;
//                That.keyOn = 2;
//                That.filter = o;
//                That.showAccountList(data);
//            }
//        })
//    }
    
    
    Account.prototype.sortByKey = function (p) {
        $.post('/agent_employee_sort', p, function (data) {
            if (!!data && data.code == 0) {
                That.page_info = data.page_info;
                That.keyOn = 2;
                That.sort = p;
                That.showAccountList(data);
            }
        })
    }

    Account.prototype.resetAccountPsw = function (psw) {
        var id = this.accountId;
        $.post('/reset_account_psw', {
            id: id,
            new_psw: psw
        }, function (data) {
            if (!!data && data.code == 0) {
                GlobalFunc.Alert.success("重设成功")
                $ModalResetPsw.modal('hide');
                $("#resetPswForm")[0].reset();
            } else {
                GlobalFunc.Alert.fail("系统出错")
            }
        })
    }

    /*---------------------------------------------------------------------------------------------------
    ---------------------------------------------------------------------------------*/

    AccountCtl.getAccountList(1);

    /*事件绑定*/
    Dom.find('.modal').modal({
        backdrop: 'static',
        keyboard: false,
        show: false
    });

    $addAccountBtn.on('click', function () {
        if (GlobalFunc.CheckForm($addAccountForm)) {
            return
        };
        var Account = new Object();
        Account.AccountId = Dom.find('#UserName').val();
        Account.Pwd = Dom.find('#Password').val();
        Account.Name = Dom.find('#Name').val();
        Account.Gender = Dom.find('#Sex').find('option:checked').val();
        Account.Mail = Dom.find('#ContactsMail').val();
        Account.Mobile = Dom.find('#ContactsPhone').val();

        AccountCtl.addAccount(Account);
    });

    $searchAccountBtn.on('click', function () {
        AccountCtl.keyOn = true;
        var key = Dom.find('#keyWord').val();
        if (!$.trim(key)) {
            AccountCtl.getAccountList(1);
            return
        }
        AccountCtl.searchAccount(key)
    })

    $delAccountBtn.on('click', function () {
        var chk_value = [];
        $('input[name="DelEmployeeCheck"]:checked').each(function () {
            chk_value.push($(this).val());
        });
        AccountCtl.deleteAccount(chk_value)
    });

    $delAccountNow.on('click', function () {
        !GlobalFunc.CheckChoose('DelEmployeeCheck') ? $ModalDelAccount.modal('show') : GlobalFunc.Alert.info('请勾选你需要删除的账号！');
    });

    $changeAccountBtn.on('click', function () {
        if (GlobalFunc.CheckForm($changeAccountForm)) {
            return
        };
        AccountCtl.changeAccount();
        AccountCtl.updateAccount(AccountCtl.Account);
    });

    Dom.on('click', '.changeAccountBtn', function () {
        $ModalChangeAccount.modal('show');
        var parent_dom = $(this).parents('tr');
        var Account = new Object();
        Account.accountId = parent_dom.find('.account_id').text();
        Account.name = parent_dom.find('.name').text();
        Account.phone = parent_dom.find('.phone').text();
        Account.mail = parent_dom.find('.mail').text();
        Account.sex = parent_dom.find('.sex').text();

        Dom.find('#_UserName').val(Account.accountId);
        Dom.find('#_Name').val(Account.name);
        Dom.find('#_ContactsPhone').val(Account.phone);
        Dom.find('#_ContactsMail').val(Account.mail);

        Dom.find('#_Sex').find('option').each(function () {
            $(this).text() == Account.sex && ($(this)[0].selected = true);
        });
    });

    $addAccountNow.click(function () {
        $ModalAddAccount.modal('show');
    });

    Dom.on('click', '.updatePsw', function () {
        $ModalResetPsw.modal('show');
        AccountCtl.accountId = $(this).parents('tr').find('.account_id').text();
    });

    $resetPswBtn.click(function () {
        if (GlobalFunc.CheckForm($('#resetPswForm'))) {
            return
        };
        var psw = $('#password_reset').val();
        AccountCtl.resetAccountPsw(psw);
    })
    
    /**
     * Sort Button [Onclick]
     */
    Dom.on('click', '.orderByBtn', function () {
        var $this = $(this);
        if ($this.hasClass('on')) {
            var DOM_I = $this.find('i.fa.on');
            DOM_I.removeClass('on').siblings('i').addClass('on');
        } else {
            $this.find('i.fa:eq(1)').addClass('on');
        }
        $this.siblings('.orderByBtn').removeClass('on').find('i.fa').removeClass('on');
        $this.addClass('on');
        
        var sort = {
            sortKey: $this.data('key'),
            sortIndex: $this.find('i.fa.on').index(),
            page: 1
        };
        AccountCtl.sortByKey(sort);
    });


    /*排序*/
//    Dom.on('click', '#timexSortBtn', function () {
//        var dom = $(this).parents('.panel-body').first();
//        var key = $(this).data('key');
//        var sortIndex = dom.find('.orderBy').find('input[type=radio]:checked').data('index');
//
//        AccountCtl.filter_info.sortIndex = sortIndex;
//        AccountCtl.filter_info.sortKey = key;
//        if (key == 'timex') {
//            AccountCtl.filter_info.timex_begin = dom.find('.startValue').val();
//            AccountCtl.filter_info.timex_end = dom.find('.endValue').val();
//        }
//
//        AccountCtl.accountFilter(AccountCtl.filter_info);
//    });

    /*翻页*/
    $Account_page.on('click', '.trun-left', function () {
        var page = AccountCtl.page_info.page - 1 || 1;
        AccountCtl.Flip(page);
    })
    $Account_page.on('click', '.turn-right', function () {
        var page = AccountCtl.page_info.page = AccountCtl.page_info.page + 1 > AccountCtl.page_info.page_total ? AccountCtl.page_info.page_total : AccountCtl.page_info.page + 1;
        AccountCtl.Flip(page);
    })
    $Account_page.on('click', '.branch-btn', function () {
        var page = Number($(this).text());
        AccountCtl.Flip(page);
    })


})($(document.body))