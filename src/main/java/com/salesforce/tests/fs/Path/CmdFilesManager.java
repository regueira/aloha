package com.salesforce.tests.fs.Path;

import java.util.Objects;

public class CmdFilesManager {

    public static final String DS = "/";

    private final CmdFiles root;

    private CmdFiles currDirectory;

    private static CmdFilesManager INSTANCE;

    private CmdFilesManager() {
        this.root = new CmdFiles(DS, Boolean.TRUE);
        currDirectory = root;
    }

    public static CmdFilesManager getInstance() {
        if (Objects.isNull(INSTANCE)) {
            INSTANCE = new CmdFilesManager();
        }
        return INSTANCE;
    }

    public CmdFiles getCurrDirectory() {
        return currDirectory;
    }

    public String createDirectory(final String filename) {
        return createResource(filename, Boolean.TRUE);
    }

    public String createFile(final String filename) {
        return createResource(filename, Boolean.FALSE);
    }

    private String createResource(final String filename, final Boolean dir) {
        CmdFiles curr = new CmdFiles(currDirectory, filename, dir);
        currDirectory.getChildren().put(filename, curr);
        return curr.getFullName();
    }

    public CmdFiles changeDirectory(final String arg) {

        //TODO: fullpath directory
        if (arg.equals("..")) {
            if (Objects.nonNull(currDirectory.getParent())) {
                currDirectory = currDirectory.getParent();
            } else {
                currDirectory = root;
            }
        } else {
            currDirectory = currDirectory.getChildren().values().stream()
                .filter(curr -> curr.isDirectory() && curr.getName().equals(arg))
                .findFirst()
                .orElseThrow(() -> new RuntimeException("Directory not found"));
        }
        return currDirectory;
    }
}
