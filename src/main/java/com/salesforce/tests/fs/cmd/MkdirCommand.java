package com.salesforce.tests.fs.cmd;

import com.salesforce.tests.fs.Path.CmdFilesManager;

import java.util.Objects;

public class MkdirCommand implements CmdCommand<String, String[]> {
    @Override
    public String execute(final String[] args) {
        String rta = "Invalid Directory name";

        if (Objects.nonNull(args)) {
            rta = CmdFilesManager.getInstance().createDirectory(args[0]);
        }

        return rta;
    }
}
