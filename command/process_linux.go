package command

import (
    "errors"
    "fmt"
    "os"
    "os/exec"
    "strconv"
    "strings"
    "syscall"
)

func (r *Runner) preProcess() error {
    // 1.init command pgid
    if r.command.SysProcAttr == nil {
        r.command.SysProcAttr = &syscall.SysProcAttr{
            Setpgid: true,
            Pgid: 0,
        }
    }
    // 2.init command execute Env
    var env []string
    if r.command.Env == nil || len(r.command.Env) == 0 {
        env = os.Environ()
    } else {
        env = r.command.Env
    }
    // 3.set HOME
    if r.homeDir != "" {
        homeEnv := fmt.Sprintf("HOME=%s", r.homeDir)
        env = append(env, homeEnv)
    }
    r.command.Env = env
    // 4.set user
    if r.user != "" {
        uid, gid, groups, err := getUserCredentials(r.user)
        if err != nil {
            return err
        }
        if r.command.SysProcAttr == nil {
            r.command.SysProcAttr = &syscall.SysProcAttr{}
        }
        r.command.SysProcAttr.Credential = &syscall.Credential{Uid: uid, Gid: gid, Groups: groups, NoSetGroups: false}
    }

    return nil
}

func getUserCredentials(sessionUser string) (uint32, uint32, []uint32, error) {
    uidCmdArgs := append([]string{"-c"}, fmt.Sprintf("id -u %s", sessionUser))
    cmd := exec.Command("sh", uidCmdArgs...)
    out, err := cmd.Output()
    if err != nil {
        fmt.Printf("Failed to retrieve uid for %s: %v", sessionUser, err)
        return 0, 0, nil, err
    }
    
    uid, err := strconv.Atoi(strings.TrimSpace(string(out)))
    if err != nil {
       fmt.Printf("%s not found: %v", sessionUser, err)
        return 0, 0, nil, err
    }
    
    gidCmdArgs := append([]string{"-c"}, fmt.Sprintf("id -g %s", sessionUser))
    cmd = exec.Command("sh", gidCmdArgs...)
    out, err = cmd.Output()
    if err != nil {
        fmt.Printf("Failed to retrieve gid for %s: %v", sessionUser, err)
        return 0, 0, nil, err
    }
    
    gid, err := strconv.Atoi(strings.TrimSpace(string(out)))
    if err != nil {
        fmt.Printf("%s not found: %v", sessionUser, err)
        return 0, 0, nil, err
    }
    
    // Get the list of associated groups
    groupNamesCmdArgs := append([]string{"-c"}, fmt.Sprintf("id %s", sessionUser))
    cmd = exec.Command("sh", groupNamesCmdArgs...)
    out, err = cmd.Output()
    if err != nil {
        fmt.Printf("Failed to retrieve groups for %s: %v", sessionUser, err)
        return 0, 0, nil, err
    }
    
    // Example : uid=1873601143(ssm-user) gid=1873600513(domain users) groups=1873600513(domain users),1873601620(joiners),1873601125(aws delegated add workstations to domain users)
    // Extract groups from the output
    groupsIndex := strings.Index(string(out), groupsIdentifier)
    var groupIds []uint32
    
    if groupsIndex > 0 {
        // Extract groups names and ids from the output
        groupNamesAndIds := strings.Split(string(out)[groupsIndex+len(groupsIdentifier):], ",")
        
        // Extract group ids from the output
        for _, value := range groupNamesAndIds {
            groupId, err := strconv.Atoi(strings.TrimSpace(value[:strings.Index(value, "(")]))
            if err != nil {
                fmt.Printf("Failed to retrieve group id from %s: %v", value, err)
                return 0, 0, nil, err
            }
            
            groupIds = append(groupIds, uint32(groupId))
        }
    }
    
    // Make sure they are non-zero valid positive ids
    if uid > 0 && gid > 0 {
        return uint32(uid), uint32(gid), groupIds, nil
    }
    
    return 0, 0, nil, errors.New("invalid uid and gid")
}

func (r *Runner) removeCredential () error {
    return nil
}
