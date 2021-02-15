#!/bin/bash -x


#show_sys_info - A script to produce a html file with system info

#### Constants

TITLE="System Information for $HOSTNAME"
RIGHT_NOW="$(date +"%x %r %Z")"
TIMESTAMP="Updated on $RIGHT_NOW by $USER"

#### Functions

system_info()
{
    echo "<p>Not implemented</p>"
}

show_uptime()
{
    echo "<h2>System uptime</h2>"
    echo "<pre>"
    uptime
    echo "</pre>"
}

drive_space()
{
    echo "<h2>Filesystem space</h2>"
    echo "<pre>"
    df
    echo "</pre>"  
}


home_space()
{
    echo "<h2>Home directory space by user</h2>"
    echo "<pre>"
    echo "Bytes Directory"
      # Turn tracing on
      set -x 
      du -s /home/* | sort -nr
      set +x
    echo "</pre>" 
}

write_page()
{
 cat <<- HTML
   <html>
    <head>
      <title>
      $TITLE  
      </title>
    </head>

    <body>
      <h1>$title</h1>
      <p>$TIMESTAMP</p> 
      $(system_info)
      $(show_uptime)
      $(drive_space)
      $(home_space)    
    </body>
   </html>
HTML
}

usage()
{
  echo "usage: show_sys_info [[[-f file ] [-i]] | [-h]]"]
}

#### Main

interactive=
# Set default filename
filename=~/sys_info_page.html

# Use positional parameters to read command line options - $0 is the progtam name, $1 - $9 are parameters
# $# variable contains a number of items on the commad line in addition to the name of the command.
# shift is a shell builtin that operates on the positional params. Each time we invoke shift, it "shifts" all the
# positional params down by one
while [ "$1" != "" ]; do
  case $1 in
    -f ) shift
         filename=$1
         ;;
    -i ) interactive=1
         ;;
    -h ) usage
         exit
         ;;
     * ) usage
         exit 1
  esac
  shift
done


if [ "$interactive" = "1" ]; then 
   response=
   read -p "Enter name of output file [$filename] > " response
   if [ -n "$response" ]; then
     filename="$response"
   fi
   
   if [ -f $filename ]; then
     # -n causes the echo command to keep the cursor on the same line.
     echo -n "Output file exists. Overwrite? (y/n) > "
     # read command has -p, -t, -s flags. -p to precede the input with a prompt, -t to introduce a timeout, -s to hi	 # de user input e.g. for password input
     read response
     if [ "$response" != "y" ]; then
       echo "Exiting program."
       exit 1
     fi 
   fi
fi
 
# Test code to verify command line processing

if [ "$interactive" = "1" ]; then
  echo "interactive is on"
else
  echo "interactive is off"
fi
echo "output file = $filename"


# Write page (comment out until testing is complete)

# write_page > $filename

