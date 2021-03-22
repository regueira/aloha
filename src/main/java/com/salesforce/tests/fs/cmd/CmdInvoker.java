package com.salesforce.tests.fs.cmd;

import java.util.HashMap;
import java.util.Map;
import java.util.Optional;

public class CmdInvoker {

    private final Map<String, CmdCommand<String, String[]>> commandMap = new HashMap<>();

    public CmdInvoker() {
        commandMap.put("pwd", new PwdCommand());
        commandMap.put("ls", new LsCommand());
        commandMap.put("mkdir", new MkdirCommand());
        commandMap.put("cd", new CdCommand());
        commandMap.put("touch", new TouchCommand());

    }

    public CmdCommand<String, String[]> execute(final String cmd) {
        return Optional.ofNullable(cmd)
                .map(commandMap::get)
                .orElseThrow(() -> new RuntimeException("command not found"));
    }

}
