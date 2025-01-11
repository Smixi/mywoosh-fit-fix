
## Archived, no more needed ! 
Since the december update, it seems MyWoosh fixed the average data issue and now also output laps !

# mywoosh-fit-fix
A simple go script to patch fit files from mywoosh.

Mywoosh produce a FIT files that you can download directly on their website, or on Strava if you have enabled Connections.

It seems that strava read the entire fit files to produce some values, like average power, heartrate and cadence.

The issue is that Garmin Connect try to read the Session messages, which does contains zeroed values for the average_power, etc.

This script simply read all your .fit files in the current directory, and create a patched_xxx.fit file that contains the average values.
