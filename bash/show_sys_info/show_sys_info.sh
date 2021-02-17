#!/bin/bash


#show_sys_info - A script to produce a html file with system info

#### Constants

TITLE="System Information for $HOSTNAME"
RIGHT_NOW="$(date +"%x %r %Z")"
TIMESTAMP="Updated on $RIGHT_NOW by $USER"

#### Functions

system_info()
{
    #https://www.cyberciti.biz/faq/mac-osx-find-tell-operating-system-version-from-bash-prompt/
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
    format="%8s%10s%10s %-s\n"
    printf "$format" "Dirs" "Files" "Blocks" "Directory"
    printf "$format" "----" "-----" "------" "---------"
    if [ $(id -u) = "0" ]; then    
    	dir_list="/Users/*"
    else
	dir_list=$HOME
    fi
    for home_dir in $dir_list; do
	total_dirs=$(find $home_dir -type d | wc -l)
    	total_files=$(find $home_dir -type f | wc -l)
    	total_blocks=$(du -s $home_dir)
     	printf "$format" "$total_dirs" "$total_files" "$total_blocks"
    done
   
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


write_page > $filename

