package com.salesforce.tests.fs.cmd;

import com.salesforce.tests.fs.Path.CmdFiles;
import com.salesforce.tests.fs.Path.CmdFilesManager;

import java.util.List;
import java.util.Map;
import java.util.Objects;
import java.util.stream.Collectors;

public class LsCommand implements CmdCommand<String, String[]> {

    @Override
    public String execute(final String[] args) {
        final Map<String, CmdFiles> children = CmdFilesManager.getInstance().getCurrDirectory().getChildren();
        return String.join("\n", executeRecursive(children, args));
    }


    private List<String> executeRecursive(final Map<String, CmdFiles> children, final String[] args) {

        final List<String> files = children.values().stream()
                .map(CmdFiles::getFullName).collect(Collectors.toList());

        if(Objects.nonNull(args) && args[0].equals("-r")) {
            children.values().stream()
                .filter(CmdFiles::isDirectory)
                .forEach(cmdFiles -> files.addAll(executeRecursive(cmdFiles.getChildren(), args)));
        }

        return files;
    }
}
