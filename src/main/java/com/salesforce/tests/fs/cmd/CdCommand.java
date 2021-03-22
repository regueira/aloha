package com.salesforce.tests.fs.cmd;

import com.salesforce.tests.fs.Path.CmdFilesManager;

import java.util.Objects;

public class CdCommand implements CmdCommand<String, String[]> {

    @Override
    public String execute(final String[] args) {
        String rta;

        if (Objects.nonNull(args)) {
            CmdFilesManager.getInstance().changeDirectory(args[0]);
            rta = args[0];
        } else {
            rta = "Directory not found";
        }
        return rta;
    }
}
