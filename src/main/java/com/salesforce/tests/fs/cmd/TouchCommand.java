package com.salesforce.tests.fs.cmd;

import com.salesforce.tests.fs.Path.CmdFilesManager;

import java.util.Objects;

public class TouchCommand implements CmdCommand<String, String[]> {
    @Override
    public String execute(final String[] args) {
        String rta = "Error creating file";

        if (Objects.nonNull(args)) {
            rta = CmdFilesManager.getInstance().createFile(args[0]);
        }
        return rta;
    }
}
