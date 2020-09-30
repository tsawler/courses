ALTER TABLE course_sections ADD COLUMN prof_name VARCHAR(512) default '';
ALTER TABLE course_sections ADD COLUMN prof_email VARCHAR(512) default '';
ALTER TABLE course_sections ADD COLUMN teams_link VARCHAR(1024) default '';

update course_sections c set prof_name = (select prof_name from courses where id = c.id);
update course_sections c set prof_email = (select prof_email from courses where id = c.id);
update course_sections c set teams_link = (select teams_link from courses where id = c.id);

alter table courses drop column prof_name;
alter table courses drop column prof_email;
alter table courses drop column teams_link;