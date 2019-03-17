# MySQL Database Strcture

> All tables below will be tried to be generated if not existing.

Genrally, the server needs read and write access to the tables below.

## `vplan`
> Contains information about VPlan sheets.

```sql
CREATE TABLE IF NOT EXISTS `vplan` (
    `id` int(11) NOT NULL,
    `date_edit` timestamp NOT NULL 
        DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `date_for` datetime NOT NULL,
    `block` char(1) NOT NULL,
    `header` text NOT NULL,
    `footer` text NOT NULL,
    `published` tinyint(1) NOT NULL,
    `deleted` int(1) NOT NULL DEFAULT '0',

    CONSTRAINT `fk_vplan_details_vplan` 
        FOREIGN KEY (`vplan_id`) 
        REFERENCES `vplan` (`id`) 
        ON DELETE CASCADE ON UPDATE NO ACTION
);
```

## `vplan_details`
> Contains all entries and information about them for all VPlans. 

```sql
CREATE TABLE IF NOT EXISTS `vplan_details` (
    `id` int(11) NOT NULL,
    `vplan_id` int(11) NOT NULL,
    `class` varchar(45) NOT NULL,
    `time` varchar(45) NOT NULL,
    `messures` varchar(255) NOT NULL,
    `responsible` varchar(255) NOT NULL,
    `reason` int(1) NOT NULL DEFAULT '1',
    `geteilt` int(1) NOT NULL,
    `notiz` varchar(45) NOT NULL,
    `deleted` int(1) NOT NULL DEFAULT '0',
    `selected` int(1) NOT NULL DEFAULT '0'
);
```

## Content
> Contains news ticker information.

```sql
CREATE TABLE IF NOT EXISTS `content` (
    `id` int(11) NOT NULL,
    `date` datetime NOT NULL,
    `author_id` int(11) DEFAULT NULL,
    `headline` varchar(255) NOT NULL,
    `short` text NOT NULL,
    `story` text NOT NULL,
    `published` tinyint(1) NOT NULL,
    `container_id` int(11) DEFAULT NULL
);
```

## `apitoken`
> Contains all registered user API tokens.

```sql
CREATE TABLE IF NOT EXISTS `apitoken` (
    `id` int PRIMARY KEY AUTO_INCREMENT,
    `ident` text NOT NULL,
    `token` text NOT NULL,
    `expire` timestamp NOT NULL
);
```

## `usersettings`
> Contains settings users have defined in the frontend.

```sql
CREATE TABLE IF NOT EXISTS `usersettings` (
    `id` int PRIMARY KEY AUTO_INCREMENT,
    `ident` text NOT NULL,
    `class` text NOT NULL,
    `theme` text NOT NULL,
    `edited` timestamp NOT NULL 
        DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAM
);
```

## `logins`
> Contains the user login log.

```sql
CREATE TABLE IF NOT EXISTS `logins` (
    `id` int PRIMARY KEY AUTO_INCREMENT,
    `ident` text NOT NULL,
    `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `type` int NOT NULL DEFAULT 0,
    `useragent` text NOT NULL,
    `ipaddress` text NOT NULL 
);
```