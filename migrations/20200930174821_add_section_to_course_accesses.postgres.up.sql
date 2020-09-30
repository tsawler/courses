ALTER TABLE course_accesses ADD COLUMN section_id integer;
update course_accesses c set section_id = course_id;