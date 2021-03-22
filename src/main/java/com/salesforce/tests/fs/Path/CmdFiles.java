package com.salesforce.tests.fs.Path;

import java.util.HashMap;
import java.util.Map;

public class CmdFiles {

    private final String name;

    private String fullName;

    private final Boolean directory;

    private CmdFiles parent;

    private Map<String, CmdFiles> children;

    public CmdFiles(final CmdFiles parent, final String name, final Boolean directory) {
        this(name, directory);
        this.parent = parent;
        if (parent.getFullName().equals(CmdFilesManager.DS)) {
            this.fullName = parent.getFullName() + name;
        } else {
            this.fullName = parent.getFullName() + CmdFilesManager.DS + name;
        }
    }

    public CmdFiles(final String name, final Boolean directory) {
        this.fullName = CmdFilesManager.DS + "root";
        if (name.length() > 100) {
            throw new RuntimeException("Invalid filename");
        }
        this.name = name;
        this.directory = directory;
        if (directory) {
            this.children = new HashMap<>();
        }
    }

    public Boolean isDirectory() {
        return directory;
    }

    public String getName() {
        return name;
    }

    public Map<String, CmdFiles> getChildren() {
        return children;
    }

    public String getFullName() {
        return fullName;
    }

    public CmdFiles getParent() {
        return parent;
    }
}
