create table if not exists users (
    userid int not null auto_increment primary key,
    email varchar(255) not null,
    firstname varchar(255) not null,
    lastname varchar(255) not null,
    username varchar(255) not null,
    photourl varchar(75) not null,
    passhash binary(72) not null,
    unique (username),
    unique(email)
);


create table if not exists userlog (
    logid int not null auto_increment primary key,
    userid int not null,
    timeinitiated datetime not null,
    ip varchar(20) not null
);



create table if not exists channels (
    channelid int not null auto_increment primary key,
    channelname varchar(69) not null,
    channeldescription varchar(255),
    isprivate boolean not null,
    createdat datetime not null,
    creatorid int not null,
    editedat datetime,
    unique (channelname)
);

create table if not exists userchannel (
    userchannelid int not null auto_increment primary key,
    userid int not null,
    channelid int not null
);


create table if not exists messages (
    messageid int not null auto_increment primary key,
    channelid int not null,
    messagebody varchar(255) not null,
    createdat datetime not null,
    creatorid int not null,
    editedat datetime
);