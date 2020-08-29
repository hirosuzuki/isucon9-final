CREATE INDEX train_master_dtt ON train_master (date, train_class, train_name)
# DROP INDEX train_master_dtt ON train_master;

CREATE INDEX train_timetable_master_dtts ON train_timetable_master (date, train_class, train_name, station)
# 
