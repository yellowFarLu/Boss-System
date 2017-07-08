--
-- ER/Studio 8.0 SQL Code Generation
-- Company :      abc
-- Project :      boss20170329.DM1
-- Author :       Windows ÓÃ»§
--
-- Date Created : Thursday, March 30, 2017 20:05:44
-- Target DBMS : MySQL 5.x
--

-- 
-- TABLE: t_account 
--
DROP DATABASE emoa_boss;

CREATE DATABASE IF NOT EXISTS emoa_boss;
USE emoa_boss;


CREATE TABLE t_account(
    account_id    VARCHAR(64)    NOT NULL,
    agent_id      INT            NOT NULL,
    pwd           VARCHAR(64)    NOT NULL,
    NAME          VARCHAR(64)    NOT NULL,
    gender        INT            NOT NULL,
    mobile        VARCHAR(64)    NOT NULL,
    mail          VARCHAR(64)    NOT NULL,
    enabled       INT            DEFAULT 0,
    timex         DATETIME       NOT NULL,
    PRIMARY KEY (account_id)
)ENGINE=INNODB
;



-- 
-- TABLE: t_account_customers 
--

CREATE TABLE t_account_customers(
    custom_id      INT            NOT NULL,
    account_id     VARCHAR(64)    NOT NULL,
    assign_time    DATETIME       NOT NULL,
    STATUS         INT,
    PRIMARY KEY (custom_id)
)ENGINE=INNODB
;



-- 
-- TABLE: t_account_roles 
--

CREATE TABLE t_account_roles(
    role_id       INT            NOT NULL,
    account_id    VARCHAR(64)    NOT NULL,
    PRIMARY KEY (role_id, account_id)
)ENGINE=INNODB
;



-- 
-- TABLE: t_agent 
--

CREATE TABLE t_agent(
    agent_id    INT            AUTO_INCREMENT,
    NAME        VARCHAR(64)    NOT NULL,
    contacts    VARCHAR(64)    NOT NULL,
    mobile      VARCHAR(64)    NOT NULL,
    mail        VARCHAR(64)    NOT NULL,
    enabled     INT            DEFAULT 0,
    note        TEXT,
    timex       DATETIME       NOT NULL,
    PRIMARY KEY (agent_id)
)ENGINE=INNODB
;



-- 
-- TABLE: t_biz 
--

CREATE TABLE t_biz(
    biz_id           INT            auto_increment  NOT NULL,
    account_id       VARCHAR(64)     NOT NULL,
    custom_id        INT             NOT NULL,
    title            VARCHAR(128)    NOT NULL,
    content          TEXT,
    amount           FLOAT(8, 0),
    status           INT             NOT NULL,
    timex            DATETIME        NOT NULL,
    estimate_time    DATETIME        NOT NULL,
    real_time        DATETIME,
    PRIMARY KEY (biz_id)
)ENGINE=INNODB
;



-- 
-- TABLE: t_biz_comments 
--

CREATE TABLE t_biz_comments(
    comment_id    INT            NOT NULL,
    biz_id        INT            NOT NULL,
    committer     VARCHAR(64)    NOT NULL,
    comments      TEXT           NOT NULL,
    timex         DATETIME       NOT NULL,
    PRIMARY KEY (comment_id)
)ENGINE=INNODB
;



-- 
-- TABLE: t_biz_tags 
--

CREATE TABLE t_biz_tags(
    biz_id    INT         NOT NULL,
    tag_id    INT         NOT NULL,
    timex     DATETIME    NOT NULL,
    PRIMARY KEY (biz_id, tag_id)
)ENGINE=INNODB
;



-- 
-- TABLE: t_city 
--

CREATE TABLE t_city(
    city_id          INT            AUTO_INCREMENT,
    province_id      INT            NOT NULL,
    city             VARCHAR(64)    NOT NULL,
    display_order    INT            NOT NULL,
    PRIMARY KEY (city_id)
)ENGINE=INNODB
;



-- 
-- TABLE: t_customer 
--

CREATE TABLE t_customer(
    custom_id           INT            AUTO_INCREMENT,
    name                VARCHAR(64)    NOT NULL,
    rtx_number          VARCHAR(8),
    contacts            VARCHAR(64),
    phone               VARCHAR(64),
    mobile              VARCHAR(64),
    qq                  VARCHAR(64),
    mail                CHAR(64),
    agent_id            INT            NOT NULL,
    city_id             INT,
    timex               DATETIME,
    assign_status       INT            NOT NULL,
    last_follow_time    DATETIME,
    note                TEXT,
    enabled       		INT            default 0,
    PRIMARY KEY (custom_id)
)ENGINE=INNODB
;



-- 
-- TABLE: t_customer_assign_history 
--

CREATE TABLE t_customer_assign_history(
    assign_id    INT            AUTO_INCREMENT,
    custom_id    INT            NOT NULL,
    assigner     VARCHAR(64)    NOT NULL,
    agent_id     INT            NOT NULL,
    assignee     VARCHAR(64),
    timex        DATETIME       NOT NULL,
    PRIMARY KEY (assign_id)
)ENGINE=INNODB
;



-- 
-- TABLE: t_customer_comments 
--

CREATE TABLE t_customer_comments(
    comment_id    INT            AUTO_INCREMENT,
    custom_id     INT            NOT NULL,
    committer     VARCHAR(64)    NOT NULL,
    comments      TEXT           NOT NULL,
    timex         DATETIME       NOT NULL,
    PRIMARY KEY (comment_id)
)ENGINE=INNODB
;



-- 
-- TABLE: t_customer_tags 
--

CREATE TABLE t_customer_tags(
    custom_id    INT         NOT NULL,
    tag_id       INT         NOT NULL,
    timex        DATETIME    NOT NULL,
    PRIMARY KEY (custom_id, tag_id)
)ENGINE=INNODB
;



-- 
-- TABLE: t_province 
--

CREATE TABLE t_province(
    province_id      INT            auto_increment NOT NULL,
    province         VARCHAR(64)    NOT NULL,
    display_order    INT            NOT NULL,
    PRIMARY KEY (province_id)
)ENGINE=INNODB
;



-- 
-- TABLE: t_role 
--

CREATE TABLE t_role(
    role_id    INT            NOT NULL,
    name       VARCHAR(64)    NOT NULL,
    PRIMARY KEY (role_id)
)ENGINE=INNODB
;



-- 
-- TABLE: t_tag 
--

CREATE TABLE t_tag(
    tag_id    INT            AUTO_INCREMENT,
    name      VARCHAR(64)    NOT NULL,
    type      VARCHAR(20),
    note      TEXT,
    PRIMARY KEY (tag_id)
)ENGINE=INNODB
;



-- 
-- TABLE: t_account 
--

ALTER TABLE t_account ADD CONSTRAINT Reft_agent34 
    FOREIGN KEY (agent_id)
    REFERENCES t_agent(agent_id)
;


-- 
-- TABLE: t_account_customers 
--

ALTER TABLE t_account_customers ADD CONSTRAINT Reft_customer13 
    FOREIGN KEY (custom_id)
    REFERENCES t_customer(custom_id)
;

ALTER TABLE t_account_customers ADD CONSTRAINT Reft_account14 
    FOREIGN KEY (account_id)
    REFERENCES t_account(account_id)
;


-- 
-- TABLE: t_account_roles 
--

ALTER TABLE t_account_roles ADD CONSTRAINT Reft_role22 
    FOREIGN KEY (role_id)
    REFERENCES t_role(role_id)
;

ALTER TABLE t_account_roles ADD CONSTRAINT Reft_account23 
    FOREIGN KEY (account_id)
    REFERENCES t_account(account_id)
;


-- 
-- TABLE: t_biz 
--

ALTER TABLE t_biz ADD CONSTRAINT Reft_account1 
    FOREIGN KEY (account_id)
    REFERENCES t_account(account_id)
;

ALTER TABLE t_biz ADD CONSTRAINT Reft_customer2 
    FOREIGN KEY (custom_id)
    REFERENCES t_customer(custom_id)
;


-- 
-- TABLE: t_biz_comments 
--

ALTER TABLE t_biz_comments ADD CONSTRAINT Reft_biz31 
    FOREIGN KEY (biz_id)
    REFERENCES t_biz(biz_id)
;

ALTER TABLE t_biz_comments ADD CONSTRAINT Reft_account37 
    FOREIGN KEY (committer)
    REFERENCES t_account(account_id)
;


-- 
-- TABLE: t_biz_tags 
--

ALTER TABLE t_biz_tags ADD CONSTRAINT Reft_biz26 
    FOREIGN KEY (biz_id)
    REFERENCES t_biz(biz_id)
;

ALTER TABLE t_biz_tags ADD CONSTRAINT Reft_tag27 
    FOREIGN KEY (tag_id)
    REFERENCES t_tag(tag_id)
;


-- 
-- TABLE: t_city 
--

ALTER TABLE t_city ADD CONSTRAINT Reft_province4 
    FOREIGN KEY (province_id)
    REFERENCES t_province(province_id)
;


-- 
-- TABLE: t_customer 
--

ALTER TABLE t_customer ADD CONSTRAINT Reft_agent35 
    FOREIGN KEY (agent_id)
    REFERENCES t_agent(agent_id)
;

ALTER TABLE t_customer ADD CONSTRAINT Reft_city36 
    FOREIGN KEY (city_id)
    REFERENCES t_city(city_id)
;


-- 
-- TABLE: t_customer_assign_history 
--

ALTER TABLE t_customer_assign_history ADD CONSTRAINT Reft_customer20 
    FOREIGN KEY (custom_id)
    REFERENCES t_customer(custom_id)
;

ALTER TABLE t_customer_assign_history ADD CONSTRAINT Reft_agent40 
    FOREIGN KEY (agent_id)
    REFERENCES t_agent(agent_id)
;

ALTER TABLE t_customer_assign_history ADD CONSTRAINT Reft_account41 
    FOREIGN KEY (assigner)
    REFERENCES t_account(account_id)
;

ALTER TABLE t_customer_assign_history ADD CONSTRAINT Reft_account42 
    FOREIGN KEY (assignee)
    REFERENCES t_account(account_id)
;


-- 
-- TABLE: t_customer_comments 
--

ALTER TABLE t_customer_comments ADD CONSTRAINT Reft_customer21 
    FOREIGN KEY (custom_id)
    REFERENCES t_customer(custom_id)
;

ALTER TABLE t_customer_comments ADD CONSTRAINT Reft_account38 
    FOREIGN KEY (committer)
    REFERENCES t_account(account_id)
;


-- 
-- TABLE: t_customer_tags 
--

ALTER TABLE t_customer_tags ADD CONSTRAINT Reft_customer24 
    FOREIGN KEY (custom_id)
    REFERENCES t_customer(custom_id)
;

ALTER TABLE t_customer_tags ADD CONSTRAINT Reft_tag25 
    FOREIGN KEY (tag_id)
    REFERENCES t_tag(tag_id)
;



#######################init##########################

#初始化Tag

#单位性质
insert into t_tag(name, type, note) values('单位性质', 'unit_nature', '国企'); 
insert into t_tag(name, type, note) values('单位性质', 'unit_nature', '私企');
insert into t_tag(name, type, note) values('单位性质', 'unit_nature', '股份制');
insert into t_tag(name, type, note) values('单位性质', 'unit_nature', '其他');

#行业划分
insert into t_tag(name, type, note) values('行业划分', 'bs_domain', '教育');
insert into t_tag(name, type, note) values('行业划分', 'bs_domain', '金融');
insert into t_tag(name, type, note) values('行业划分', 'bs_domain', '石油石化');
insert into t_tag(name, type, note) values('行业划分', 'bs_domain', '能源');
insert into t_tag(name, type, note) values('行业划分', 'bs_domain', '钢铁');
insert into t_tag(name, type, note) values('行业划分', 'bs_domain', '餐饮');
insert into t_tag(name, type, note) values('行业划分', 'bs_domain', '电讯业');
insert into t_tag(name, type, note) values('行业划分', 'bs_domain', '房地产');
insert into t_tag(name, type, note) values('行业划分', 'bs_domain', '媒体');
insert into t_tag(name, type, note) values('行业划分', 'bs_domain', '出版社');
insert into t_tag(name, type, note) values('行业划分', 'bs_domain', '医疗');
insert into t_tag(name, type, note) values('行业划分', 'bs_domain', '互联网');
insert into t_tag(name, type, note) values('行业划分', 'bs_domain', '体育运动');
insert into t_tag(name, type, note) values('行业划分', 'bs_domain', '政府机关');
insert into t_tag(name, type, note) values('行业划分', 'bs_domain', '其他');

#客户级别
insert into t_tag(name, type, note) values('客户级别', 'cu_level', 'A类');
insert into t_tag(name, type, note) values('客户级别', 'cu_level', 'B类');
insert into t_tag(name, type, note) values('客户级别', 'cu_level', 'C类');

#商机类别
insert into t_tag(name, type, note) values('商机类别', 'biz_category', '有度销售');
insert into t_tag(name, type, note) values('商机类别', 'biz_category', 'RTX销售');
insert into t_tag(name, type, note) values('商机类别', 'biz_category', 'RTX短信');
insert into t_tag(name, type, note) values('商机类别', 'biz_category', '企业微信项目');
insert into t_tag(name, type, note) values('商机类别', 'biz_category', '有度项目');





#省份城市
insert into t_province values(1,'北京',0);
insert into t_province values(2,'天津',0);
insert into t_province values(3,'上海',0);
insert into t_province values(4,'重庆',0);
insert into t_province values(5,'河北',0);
insert into t_province values(6,'山西',0);
insert into t_province values(7,'台湾',0);
insert into t_province values(8,'辽宁',0);
insert into t_province values(9,'吉林',0);
insert into t_province values(10,'黑龙江',0);
insert into t_province values(11,'江苏',0);
insert into t_province values(12,'浙江',0);
insert into t_province values(13,'安徽',0);
insert into t_province values(14,'福建',0);
insert into t_province values(15,'江西',0);
insert into t_province values(16,'山东',0);
insert into t_province values(17,'河南',0);
insert into t_province values(18,'湖北',0);
insert into t_province values(19,'湖南',0);
insert into t_province values(20,'广东',0);
insert into t_province values(21,'甘肃',0);
insert into t_province values(22,'四川',0);
insert into t_province values(24,'贵州',0);
insert into t_province values(25,'海南',0);
insert into t_province values(26,'云南',0);
insert into t_province values(27,'青海',0);
insert into t_province values(28,'陕西',0);
insert into t_province values(29,'广西',0);
insert into t_province values(30,'西藏',0);
insert into t_province values(31,'宁夏',0);
insert into t_province values(32,'新疆',0);
insert into t_province values(33,'内蒙古',0);
insert into t_province values(34,'澳门',0);
insert into t_province values(35,'香港',0);



insert into t_city(city, province_id, display_order) values('北京',1,0);
insert into t_city(city, province_id, display_order) values('天津',2,0);
insert into t_city(city, province_id, display_order) values('上海',3,0);
insert into t_city(city, province_id, display_order) values('重庆',4,0);
insert into t_city(city, province_id, display_order) values('石家庄',5,0);
insert into t_city(city, province_id, display_order) values('唐山',5,0);
insert into t_city(city, province_id, display_order) values('秦皇岛',5,0);
insert into t_city(city, province_id, display_order) values('邯郸',5,0);
insert into t_city(city, province_id, display_order) values('邢台',5,0);
insert into t_city(city, province_id, display_order) values('保定',5,0);
insert into t_city(city, province_id, display_order) values('张家口',5,0);
insert into t_city(city, province_id, display_order) values('承德',5,0);
insert into t_city(city, province_id, display_order) values('沧州',5,0);
insert into t_city(city, province_id, display_order) values('廊坊',5,0);
insert into t_city(city, province_id, display_order) values('衡水',5,0);
insert into t_city(city, province_id, display_order) values('太原',6,0);
insert into t_city(city, province_id, display_order) values('大同',6,0);
insert into t_city(city, province_id, display_order) values('阳泉',6,0);
insert into t_city(city, province_id, display_order) values('长治',6,0);
insert into t_city(city, province_id, display_order) values('晋城',6,0);
insert into t_city(city, province_id, display_order) values('朔州',6,0);
insert into t_city(city, province_id, display_order) values('晋中',6,0);
insert into t_city(city, province_id, display_order) values('运城',6,0);
insert into t_city(city, province_id, display_order) values('忻州',6,0);
insert into t_city(city, province_id, display_order) values('临汾',6,0);
insert into t_city(city, province_id, display_order) values('吕梁',6,0);
insert into t_city(city, province_id, display_order) values('台北',7,0);
insert into t_city(city, province_id, display_order) values('高雄',7,0);
insert into t_city(city, province_id, display_order) values('基隆',7,0);
insert into t_city(city, province_id, display_order) values('台中',7,0);
insert into t_city(city, province_id, display_order) values('台南',7,0);
insert into t_city(city, province_id, display_order) values('新竹',7,0);
insert into t_city(city, province_id, display_order) values('嘉义',7,0);
insert into t_city(city, province_id, display_order) values('台北',7,0);
insert into t_city(city, province_id, display_order) values('宜兰',7,0);
insert into t_city(city, province_id, display_order) values('桃园',7,0);
insert into t_city(city, province_id, display_order) values('新竹',7,0);
insert into t_city(city, province_id, display_order) values('苗栗',7,0);
insert into t_city(city, province_id, display_order) values('台中',7,0);
insert into t_city(city, province_id, display_order) values('彰化',7,0);
insert into t_city(city, province_id, display_order) values('南投',7,0);
insert into t_city(city, province_id, display_order) values('云林',7,0);
insert into t_city(city, province_id, display_order) values('嘉义',7,0);
insert into t_city(city, province_id, display_order) values('台南',7,0);
insert into t_city(city, province_id, display_order) values('高雄',7,0);
insert into t_city(city, province_id, display_order) values('屏东',7,0);
insert into t_city(city, province_id, display_order) values('澎湖',7,0);
insert into t_city(city, province_id, display_order) values('台东',7,0);
insert into t_city(city, province_id, display_order) values('花莲',7,0);
insert into t_city(city, province_id, display_order) values('沈阳',8,0);
insert into t_city(city, province_id, display_order) values('大连',8,0);
insert into t_city(city, province_id, display_order) values('鞍山',8,0);
insert into t_city(city, province_id, display_order) values('抚顺',8,0);
insert into t_city(city, province_id, display_order) values('本溪',8,0);
insert into t_city(city, province_id, display_order) values('丹东',8,0);
insert into t_city(city, province_id, display_order) values('锦州',8,0);
insert into t_city(city, province_id, display_order) values('营口',8,0);
insert into t_city(city, province_id, display_order) values('阜新',8,0);
insert into t_city(city, province_id, display_order) values('辽阳',8,0);
insert into t_city(city, province_id, display_order) values('盘锦',8,0);
insert into t_city(city, province_id, display_order) values('铁岭',8,0);
insert into t_city(city, province_id, display_order) values('朝阳',8,0);
insert into t_city(city, province_id, display_order) values('葫芦岛',8,0);
insert into t_city(city, province_id, display_order) values('长春',9,0);
insert into t_city(city, province_id, display_order) values('吉林',9,0);
insert into t_city(city, province_id, display_order) values('四平',9,0);
insert into t_city(city, province_id, display_order) values('辽源',9,0);
insert into t_city(city, province_id, display_order) values('通化',9,0);
insert into t_city(city, province_id, display_order) values('白山',9,0);
insert into t_city(city, province_id, display_order) values('松原',9,0);
insert into t_city(city, province_id, display_order) values('白城',9,0);
insert into t_city(city, province_id, display_order) values('延边',9,0);
insert into t_city(city, province_id, display_order) values('哈尔滨',10,0);
insert into t_city(city, province_id, display_order) values('齐齐哈尔',10,0);
insert into t_city(city, province_id, display_order) values('鹤岗',10,0);
insert into t_city(city, province_id, display_order) values('双鸭山',10,0);
insert into t_city(city, province_id, display_order) values('鸡西',10,0);
insert into t_city(city, province_id, display_order) values('大庆',10,0);
insert into t_city(city, province_id, display_order) values('伊春',10,0);
insert into t_city(city, province_id, display_order) values('牡丹江',10,0);
insert into t_city(city, province_id, display_order) values('佳木斯',10,0);
insert into t_city(city, province_id, display_order) values('七台河',10,0);
insert into t_city(city, province_id, display_order) values('黑河',10,0);
insert into t_city(city, province_id, display_order) values('绥化',10,0);
insert into t_city(city, province_id, display_order) values('大兴安岭',10,0);
insert into t_city(city, province_id, display_order) values('南京',11,0);
insert into t_city(city, province_id, display_order) values('无锡',11,0);
insert into t_city(city, province_id, display_order) values('徐州',11,0);
insert into t_city(city, province_id, display_order) values('常州',11,0);
insert into t_city(city, province_id, display_order) values('苏州',11,0);
insert into t_city(city, province_id, display_order) values('南通',11,0);
insert into t_city(city, province_id, display_order) values('连云港',11,0);
insert into t_city(city, province_id, display_order) values('淮安',11,0);
insert into t_city(city, province_id, display_order) values('盐城',11,0);
insert into t_city(city, province_id, display_order) values('扬州',11,0);
insert into t_city(city, province_id, display_order) values('镇江',11,0);
insert into t_city(city, province_id, display_order) values('泰州',11,0);
insert into t_city(city, province_id, display_order) values('宿迁',11,0);
insert into t_city(city, province_id, display_order) values('杭州',12,0);
insert into t_city(city, province_id, display_order) values('宁波',12,0);
insert into t_city(city, province_id, display_order) values('温州',12,0);
insert into t_city(city, province_id, display_order) values('嘉兴',12,0);
insert into t_city(city, province_id, display_order) values('湖州',12,0);
insert into t_city(city, province_id, display_order) values('绍兴',12,0);
insert into t_city(city, province_id, display_order) values('金华',12,0);
insert into t_city(city, province_id, display_order) values('衢州',12,0);
insert into t_city(city, province_id, display_order) values('舟山',12,0);
insert into t_city(city, province_id, display_order) values('台州',12,0);
insert into t_city(city, province_id, display_order) values('丽水',12,0);
insert into t_city(city, province_id, display_order) values('合肥',13,0);
insert into t_city(city, province_id, display_order) values('芜湖',13,0);
insert into t_city(city, province_id, display_order) values('蚌埠',13,0);
insert into t_city(city, province_id, display_order) values('淮南',13,0);
insert into t_city(city, province_id, display_order) values('马鞍山',13,0);
insert into t_city(city, province_id, display_order) values('淮北',13,0);
insert into t_city(city, province_id, display_order) values('铜陵',13,0);
insert into t_city(city, province_id, display_order) values('安庆',13,0);
insert into t_city(city, province_id, display_order) values('黄山',13,0);
insert into t_city(city, province_id, display_order) values('滁州',13,0);
insert into t_city(city, province_id, display_order) values('阜阳',13,0);
insert into t_city(city, province_id, display_order) values('宿州',13,0);
insert into t_city(city, province_id, display_order) values('巢湖',13,0);
insert into t_city(city, province_id, display_order) values('六安',13,0);
insert into t_city(city, province_id, display_order) values('亳州',13,0);
insert into t_city(city, province_id, display_order) values('池州',13,0);
insert into t_city(city, province_id, display_order) values('宣城',13,0);
insert into t_city(city, province_id, display_order) values('福州',14,0);
insert into t_city(city, province_id, display_order) values('厦门',14,0);
insert into t_city(city, province_id, display_order) values('莆田',14,0);
insert into t_city(city, province_id, display_order) values('三明',14,0);
insert into t_city(city, province_id, display_order) values('泉州',14,0);
insert into t_city(city, province_id, display_order) values('漳州',14,0);
insert into t_city(city, province_id, display_order) values('南平',14,0);
insert into t_city(city, province_id, display_order) values('龙岩',14,0);
insert into t_city(city, province_id, display_order) values('宁德',14,0);
insert into t_city(city, province_id, display_order) values('南昌',15,0);
insert into t_city(city, province_id, display_order) values('景德镇',15,0);
insert into t_city(city, province_id, display_order) values('萍乡',15,0);
insert into t_city(city, province_id, display_order) values('九江',15,0);
insert into t_city(city, province_id, display_order) values('新余',15,0);
insert into t_city(city, province_id, display_order) values('鹰潭',15,0);
insert into t_city(city, province_id, display_order) values('赣州',15,0);
insert into t_city(city, province_id, display_order) values('吉安',15,0);
insert into t_city(city, province_id, display_order) values('宜春',15,0);
insert into t_city(city, province_id, display_order) values('抚州',15,0);
insert into t_city(city, province_id, display_order) values('上饶',15,0);
insert into t_city(city, province_id, display_order) values('济南',16,0);
insert into t_city(city, province_id, display_order) values('青岛',16,0);
insert into t_city(city, province_id, display_order) values('淄博',16,0);
insert into t_city(city, province_id, display_order) values('枣庄',16,0);
insert into t_city(city, province_id, display_order) values('东营',16,0);
insert into t_city(city, province_id, display_order) values('烟台',16,0);
insert into t_city(city, province_id, display_order) values('潍坊',16,0);
insert into t_city(city, province_id, display_order) values('济宁',16,0);
insert into t_city(city, province_id, display_order) values('泰安',16,0);
insert into t_city(city, province_id, display_order) values('威海',16,0);
insert into t_city(city, province_id, display_order) values('日照',16,0);
insert into t_city(city, province_id, display_order) values('莱芜',16,0);
insert into t_city(city, province_id, display_order) values('临沂',16,0);
insert into t_city(city, province_id, display_order) values('德州',16,0);
insert into t_city(city, province_id, display_order) values('聊城',16,0);
insert into t_city(city, province_id, display_order) values('滨州',16,0);
insert into t_city(city, province_id, display_order) values('菏泽',16,0);
insert into t_city(city, province_id, display_order) values('郑州',17,0);
insert into t_city(city, province_id, display_order) values('开封',17,0);
insert into t_city(city, province_id, display_order) values('洛阳',17,0);
insert into t_city(city, province_id, display_order) values('平顶山',17,0);
insert into t_city(city, province_id, display_order) values('安阳',17,0);
insert into t_city(city, province_id, display_order) values('鹤壁',17,0);
insert into t_city(city, province_id, display_order) values('新乡',17,0);
insert into t_city(city, province_id, display_order) values('焦作',17,0);
insert into t_city(city, province_id, display_order) values('濮阳',17,0);
insert into t_city(city, province_id, display_order) values('许昌',17,0);
insert into t_city(city, province_id, display_order) values('漯河',17,0);
insert into t_city(city, province_id, display_order) values('三门峡',17,0);
insert into t_city(city, province_id, display_order) values('南阳',17,0);
insert into t_city(city, province_id, display_order) values('商丘',17,0);
insert into t_city(city, province_id, display_order) values('信阳',17,0);
insert into t_city(city, province_id, display_order) values('周口',17,0);
insert into t_city(city, province_id, display_order) values('驻马店',17,0);
insert into t_city(city, province_id, display_order) values('济源',17,0);
insert into t_city(city, province_id, display_order) values('武汉',18,0);
insert into t_city(city, province_id, display_order) values('黄石',18,0);
insert into t_city(city, province_id, display_order) values('十堰',18,0);
insert into t_city(city, province_id, display_order) values('荆州',18,0);
insert into t_city(city, province_id, display_order) values('宜昌',18,0);
insert into t_city(city, province_id, display_order) values('襄樊',18,0);
insert into t_city(city, province_id, display_order) values('鄂州',18,0);
insert into t_city(city, province_id, display_order) values('荆门',18,0);
insert into t_city(city, province_id, display_order) values('孝感',18,0);
insert into t_city(city, province_id, display_order) values('黄冈',18,0);
insert into t_city(city, province_id, display_order) values('咸宁',18,0);
insert into t_city(city, province_id, display_order) values('随州',18,0);
insert into t_city(city, province_id, display_order) values('仙桃',18,0);
insert into t_city(city, province_id, display_order) values('天门',18,0);
insert into t_city(city, province_id, display_order) values('潜江',18,0);
insert into t_city(city, province_id, display_order) values('神农架',18,0);
insert into t_city(city, province_id, display_order) values('恩施',18,0);
insert into t_city(city, province_id, display_order) values('长沙',19,0);
insert into t_city(city, province_id, display_order) values('株洲',19,0);
insert into t_city(city, province_id, display_order) values('湘潭',19,0);
insert into t_city(city, province_id, display_order) values('衡阳',19,0);
insert into t_city(city, province_id, display_order) values('邵阳',19,0);
insert into t_city(city, province_id, display_order) values('岳阳',19,0);
insert into t_city(city, province_id, display_order) values('常德',19,0);
insert into t_city(city, province_id, display_order) values('张家界',19,0);
insert into t_city(city, province_id, display_order) values('益阳',19,0);
insert into t_city(city, province_id, display_order) values('郴州',19,0);
insert into t_city(city, province_id, display_order) values('永州',19,0);
insert into t_city(city, province_id, display_order) values('怀化',19,0);
insert into t_city(city, province_id, display_order) values('娄底',19,0);
insert into t_city(city, province_id, display_order) values('湘西',19,0);
insert into t_city(city, province_id, display_order) values('广州',20,0);
insert into t_city(city, province_id, display_order) values('深圳',20,0);
insert into t_city(city, province_id, display_order) values('珠海',20,0);
insert into t_city(city, province_id, display_order) values('汕头',20,0);
insert into t_city(city, province_id, display_order) values('韶关',20,0);
insert into t_city(city, province_id, display_order) values('佛山',20,0);
insert into t_city(city, province_id, display_order) values('江门',20,0);
insert into t_city(city, province_id, display_order) values('湛江',20,0);
insert into t_city(city, province_id, display_order) values('茂名',20,0);
insert into t_city(city, province_id, display_order) values('肇庆',20,0);
insert into t_city(city, province_id, display_order) values('惠州',20,0);
insert into t_city(city, province_id, display_order) values('梅州',20,0);
insert into t_city(city, province_id, display_order) values('汕尾',20,0);
insert into t_city(city, province_id, display_order) values('河源',20,0);
insert into t_city(city, province_id, display_order) values('阳江',20,0);
insert into t_city(city, province_id, display_order) values('清远',20,0);
insert into t_city(city, province_id, display_order) values('东莞',20,0);
insert into t_city(city, province_id, display_order) values('中山',20,0);
insert into t_city(city, province_id, display_order) values('潮州',20,0);
insert into t_city(city, province_id, display_order) values('揭阳',20,0);
insert into t_city(city, province_id, display_order) values('云浮',20,0);
insert into t_city(city, province_id, display_order) values('兰州',21,0);
insert into t_city(city, province_id, display_order) values('金昌',21,0);
insert into t_city(city, province_id, display_order) values('白银',21,0);
insert into t_city(city, province_id, display_order) values('天水',21,0);
insert into t_city(city, province_id, display_order) values('嘉峪关',21,0);
insert into t_city(city, province_id, display_order) values('武威',21,0);
insert into t_city(city, province_id, display_order) values('张掖',21,0);
insert into t_city(city, province_id, display_order) values('平凉',21,0);
insert into t_city(city, province_id, display_order) values('酒泉',21,0);
insert into t_city(city, province_id, display_order) values('庆阳',21,0);
insert into t_city(city, province_id, display_order) values('定西',21,0);
insert into t_city(city, province_id, display_order) values('陇南',21,0);
insert into t_city(city, province_id, display_order) values('临夏',21,0);
insert into t_city(city, province_id, display_order) values('甘南',21,0);
insert into t_city(city, province_id, display_order) values('成都',22,0);
insert into t_city(city, province_id, display_order) values('自贡',22,0);
insert into t_city(city, province_id, display_order) values('攀枝花',22,0);
insert into t_city(city, province_id, display_order) values('泸州',22,0);
insert into t_city(city, province_id, display_order) values('德阳',22,0);
insert into t_city(city, province_id, display_order) values('绵阳',22,0);
insert into t_city(city, province_id, display_order) values('广元',22,0);
insert into t_city(city, province_id, display_order) values('遂宁',22,0);
insert into t_city(city, province_id, display_order) values('内江',22,0);
insert into t_city(city, province_id, display_order) values('乐山',22,0);
insert into t_city(city, province_id, display_order) values('南充',22,0);
insert into t_city(city, province_id, display_order) values('眉山',22,0);
insert into t_city(city, province_id, display_order) values('宜宾',22,0);
insert into t_city(city, province_id, display_order) values('广安',22,0);
insert into t_city(city, province_id, display_order) values('达州',22,0);
insert into t_city(city, province_id, display_order) values('雅安',22,0);
insert into t_city(city, province_id, display_order) values('巴中',22,0);
insert into t_city(city, province_id, display_order) values('资阳',22,0);
insert into t_city(city, province_id, display_order) values('阿坝',22,0);
insert into t_city(city, province_id, display_order) values('甘孜',22,0);
insert into t_city(city, province_id, display_order) values('凉山',22,0);
insert into t_city(city, province_id, display_order) values('贵阳',24,0);
insert into t_city(city, province_id, display_order) values('六盘水',24,0);
insert into t_city(city, province_id, display_order) values('遵义',24,0);
insert into t_city(city, province_id, display_order) values('安顺',24,0);
insert into t_city(city, province_id, display_order) values('铜仁',24,0);
insert into t_city(city, province_id, display_order) values('毕节',24,0);
insert into t_city(city, province_id, display_order) values('黔西南',24,0);
insert into t_city(city, province_id, display_order) values('黔东南',24,0);
insert into t_city(city, province_id, display_order) values('黔南',24,0);
insert into t_city(city, province_id, display_order) values('海口',25,0);
insert into t_city(city, province_id, display_order) values('三亚',25,0);
insert into t_city(city, province_id, display_order) values('五指山',25,0);
insert into t_city(city, province_id, display_order) values('琼海',25,0);
insert into t_city(city, province_id, display_order) values('儋州',25,0);
insert into t_city(city, province_id, display_order) values('文昌',25,0);
insert into t_city(city, province_id, display_order) values('万宁',25,0);
insert into t_city(city, province_id, display_order) values('东方',25,0);
insert into t_city(city, province_id, display_order) values('澄迈',25,0);
insert into t_city(city, province_id, display_order) values('定安',25,0);
insert into t_city(city, province_id, display_order) values('屯昌',25,0);
insert into t_city(city, province_id, display_order) values('临高',25,0);
insert into t_city(city, province_id, display_order) values('白沙',25,0);
insert into t_city(city, province_id, display_order) values('昌江',25,0);
insert into t_city(city, province_id, display_order) values('乐东',25,0);
insert into t_city(city, province_id, display_order) values('陵水',25,0);
insert into t_city(city, province_id, display_order) values('保亭',25,0);
insert into t_city(city, province_id, display_order) values('琼中',25,0);
insert into t_city(city, province_id, display_order) values('昆明',26,0);
insert into t_city(city, province_id, display_order) values('曲靖',26,0);
insert into t_city(city, province_id, display_order) values('玉溪',26,0);
insert into t_city(city, province_id, display_order) values('保山',26,0);
insert into t_city(city, province_id, display_order) values('昭通',26,0);
insert into t_city(city, province_id, display_order) values('丽江',26,0);
insert into t_city(city, province_id, display_order) values('思茅',26,0);
insert into t_city(city, province_id, display_order) values('临沧',26,0);
insert into t_city(city, province_id, display_order) values('文山',26,0);
insert into t_city(city, province_id, display_order) values('红河',26,0);
insert into t_city(city, province_id, display_order) values('西双版纳',26,0);
insert into t_city(city, province_id, display_order) values('楚雄',26,0);
insert into t_city(city, province_id, display_order) values('大理',26,0);
insert into t_city(city, province_id, display_order) values('德宏',26,0);
insert into t_city(city, province_id, display_order) values('怒江',26,0);
insert into t_city(city, province_id, display_order) values('迪庆',26,0);
insert into t_city(city, province_id, display_order) values('西宁',27,0);
insert into t_city(city, province_id, display_order) values('海东',27,0);
insert into t_city(city, province_id, display_order) values('海北',27,0);
insert into t_city(city, province_id, display_order) values('黄南',27,0);
insert into t_city(city, province_id, display_order) values('海南',27,0);
insert into t_city(city, province_id, display_order) values('果洛',27,0);
insert into t_city(city, province_id, display_order) values('玉树',27,0);
insert into t_city(city, province_id, display_order) values('海西',27,0);
insert into t_city(city, province_id, display_order) values('西安',28,0);
insert into t_city(city, province_id, display_order) values('铜川',28,0);
insert into t_city(city, province_id, display_order) values('宝鸡',28,0);
insert into t_city(city, province_id, display_order) values('咸阳',28,0);
insert into t_city(city, province_id, display_order) values('渭南',28,0);
insert into t_city(city, province_id, display_order) values('延安',28,0);
insert into t_city(city, province_id, display_order) values('汉中',28,0);
insert into t_city(city, province_id, display_order) values('榆林',28,0);
insert into t_city(city, province_id, display_order) values('安康',28,0);
insert into t_city(city, province_id, display_order) values('商洛',28,0);
insert into t_city(city, province_id, display_order) values('南宁',29,0);
insert into t_city(city, province_id, display_order) values('柳州',29,0);
insert into t_city(city, province_id, display_order) values('桂林',29,0);
insert into t_city(city, province_id, display_order) values('北海',29,0);
insert into t_city(city, province_id, display_order) values('防城港',29,0);
insert into t_city(city, province_id, display_order) values('钦州',29,0);
insert into t_city(city, province_id, display_order) values('贵港',29,0);
insert into t_city(city, province_id, display_order) values('玉林',29,0);
insert into t_city(city, province_id, display_order) values('百色',29,0);
insert into t_city(city, province_id, display_order) values('贺州',29,0);
insert into t_city(city, province_id, display_order) values('河池',29,0);
insert into t_city(city, province_id, display_order) values('来宾',29,0);
insert into t_city(city, province_id, display_order) values('崇左',29,0);
insert into t_city(city, province_id, display_order) values('拉萨',30,0);
insert into t_city(city, province_id, display_order) values('那曲',30,0);
insert into t_city(city, province_id, display_order) values('昌都',30,0);
insert into t_city(city, province_id, display_order) values('山南',30,0);
insert into t_city(city, province_id, display_order) values('日喀则',30,0);
insert into t_city(city, province_id, display_order) values('阿里',30,0);
insert into t_city(city, province_id, display_order) values('林芝',30,0);
insert into t_city(city, province_id, display_order) values('银川',31,0);
insert into t_city(city, province_id, display_order) values('石嘴山',31,0);
insert into t_city(city, province_id, display_order) values('吴忠',31,0);
insert into t_city(city, province_id, display_order) values('固原',31,0);
insert into t_city(city, province_id, display_order) values('中卫',31,0);
insert into t_city(city, province_id, display_order) values('乌鲁木齐',32,0);
insert into t_city(city, province_id, display_order) values('克拉玛依',32,0);
insert into t_city(city, province_id, display_order) values('石河子　',32,0);
insert into t_city(city, province_id, display_order) values('阿拉尔',32,0);
insert into t_city(city, province_id, display_order) values('图木舒克',32,0);
insert into t_city(city, province_id, display_order) values('五家渠',32,0);
insert into t_city(city, province_id, display_order) values('吐鲁番',32,0);
insert into t_city(city, province_id, display_order) values('阿克苏',32,0);
insert into t_city(city, province_id, display_order) values('喀什',32,0);
insert into t_city(city, province_id, display_order) values('哈密',32,0);
insert into t_city(city, province_id, display_order) values('和田',32,0);
insert into t_city(city, province_id, display_order) values('阿图什',32,0);
insert into t_city(city, province_id, display_order) values('库尔勒',32,0);
insert into t_city(city, province_id, display_order) values('昌吉　',32,0);
insert into t_city(city, province_id, display_order) values('阜康',32,0);
insert into t_city(city, province_id, display_order) values('米泉',32,0);
insert into t_city(city, province_id, display_order) values('博乐',32,0);
insert into t_city(city, province_id, display_order) values('伊宁',32,0);
insert into t_city(city, province_id, display_order) values('奎屯',32,0);
insert into t_city(city, province_id, display_order) values('塔城',32,0);
insert into t_city(city, province_id, display_order) values('乌苏',32,0);
insert into t_city(city, province_id, display_order) values('阿勒泰',32,0);
insert into t_city(city, province_id, display_order) values('呼和浩特',33,0);
insert into t_city(city, province_id, display_order) values('包头',33,0);
insert into t_city(city, province_id, display_order) values('乌海',33,0);
insert into t_city(city, province_id, display_ort_accountder) values('赤峰',33,0);
insert into t_city(city, province_id, display_order) values('通辽',33,0);
insert into t_city(city, province_id, display_order) values('鄂尔多斯',33,0);
insert into t_city(city, province_id, display_order) values('呼伦贝尔',33,0);
insert into t_city(city, province_id, display_order) values('巴彦淖尔',33,0);
insert into t_city(city, province_id, display_order) values('乌兰察布',33,0);
insert into t_city(city, province_id, display_order) values('锡林郭勒盟',33,0);
insert into t_city(city, province_id, display_order) values('兴安盟',33,0);
insert into t_city(city, province_id, display_order) values('阿拉善盟',33,0);
insert into t_city(city, province_id, display_order) values('澳门',34,0);
insert into t_city(city, province_id, display_order) values('香港',35,0);


INSERT INTO t_role VALUES (0, '超级管理员');
INSERT INTO t_role VALUES (1, '信达管理员');
INSERT INTO t_role VALUES (2, '渠道');
INSERT INTO t_role VALUES (3, '员工');


#信达渠道
INSERT INTO t_agent(name, contacts, mobile, mail, note, timex) VALUES ('未分配', 'wu', '13527206719', '296947440@qq.com', '广东 珠海', '2017-03-31');
INSERT INTO t_agent(name, contacts, mobile, mail, note, timex) VALUES ('信达九州', 'wu', '13527206719', '296947440@qq.com', '广东 珠海', '2017-03-31');
INSERT INTO t_account VALUES ('admin', 2,  '9806fc19405e136195fab8da8e2f922faf019430', 'wu', 0, 5351124,  '771412508@qq.com', 0, '2017-03-31');
insert into t_account_roles values(0, 'admin');

######################
select t1.account_id,t1.agent_id, t1.gender, t1.mobile, t1.name, 
	t1.pwd, t1.timex, t2.role_id from t_account t1 LEFT JOIN t_account_roles t2 ON t1.account_id=t2.account_id
	where t1.account_id = 'peyton.li@xinda.im';