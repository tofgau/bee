# bee - A simple - monolitic & autonomous - system to monitor, keep tracks and alarms what you want. Just upload an HTTPx state.
If have been working more than ten years as subcontractor in places where nothing is working easily.

I am recoding in GO "bee", a tool I developped ten years ago in oder to empower enginners.
Engineers are able to check theirs platform. They know what to do or what to run in order to say if yes or no, a system is working good.
So why is it always so complex to deal with monitoring ?

The idea with bee is to tell all people "well run the check/the job you want, run it where you want, just tell me result using this rest API"

Bee will provide 
* a simple view of all running checks/running jobs  and who (ip/path) provide them 
* a management of "time" - I want to be sure all my check are running
* a mangement of "flapping" - we we get ok, not ok ,ok , not ok ...
* simple integration with your alarming system

**Visit the wiki for the details**

# How to send a state with powershell
Example here use http, https need a valid ssl certificate but works the same way (or try to alter the invoke-webrequest command to skip it but this is really not recommanded)

$Parameters = @{

  UID = 'TEST.test'
  AuthKey='Whatever you want the first time, the same after'
  
  DOC_revision='2.1'
  
  S_currentState='OK'
  
  RunningPath=get-location
  
  RunningLocation=hostname
}

 Invoke-WebRequest -Uri 'http://<ip:port>/update' -Body $Parameters -Method Post
 
 
< UNDER CONSTRUCTION >
