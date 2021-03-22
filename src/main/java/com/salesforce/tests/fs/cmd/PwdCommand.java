package com.salesforce.tests.fs.cmd;

import com.salesforce.tests.fs.Path.CmdFilesManager;

public class PwdCommand implements CmdCommand<String, String[]> {

    @Override
    public String execute(final String[] args) {
        return CmdFilesManager.getInstance().getCurrDirectory().getFullName();
    }
}
