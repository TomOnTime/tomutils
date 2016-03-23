alias gp='cd ~/gitwork/puppet'
alias gph='cd ~/gitwork/puppet/hieradata'
alias gphh='cd ~/gitwork/puppet/hieradata/host'

# "gpm" go to the modules directory.
# "gpm foo" go to the "manifests" subdirectory of module foo.
#            If foo doesn't exist, look for modules that start with f-o-o.
function gpm() {
   cd ~/gitwork/puppet/modules
   if [[ ! -z "$1" ]]; then
    if [[ -d "$1/manifests" ]]; then
      cd "$1/manifests"
    else
      local l=( $1*/manifests )
      echo FOUND=$l
      cd $l
   fi
 fi
}
