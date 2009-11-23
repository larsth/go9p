// Copyright 2009 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package clnt

import "plan9/p"
import "syscall"

// Creates an authentication fid for the specified user. Returns the fid, if
// successful, or an Error.
func (clnt *Clnt) Auth(user p.User, aname string) (*Fid, *p.Error) {
	fid := clnt.FidAlloc();
	tc := p.NewFcall(clnt.Msize);
	err := p.PackTauth(tc, fid.Fid, user.Name(), aname, uint32(user.Id()), clnt.Dotu);
	if err != nil {
		return nil, err
	}

	_, err = clnt.rpc(tc);
	if err != nil {
		return nil, err
	}

	fid.User = user;
	return fid, nil;
}

// Creates a fid for the specified user that points to the root
// of the file server's file tree. Returns a Fid pointing to the root,
// if successful, or an Error.
func (clnt *Clnt) Attach(afid *Fid, user p.User, aname string) (*Fid, *p.Error) {
	var afno uint32;

	if afid != nil {
		afno = afid.Fid
	} else {
		afno = p.Nofid
	}

	fid := clnt.FidAlloc();
	tc := p.NewFcall(clnt.Msize);
	err := p.PackTattach(tc, fid.Fid, afno, user.Name(), aname, uint32(user.Id()), clnt.Dotu);
	if err != nil {
		return nil, err
	}

	rc, err := clnt.rpc(tc);
	if err != nil {
		return nil, err
	}

	fid.Fqid = rc.Fqid;
	fid.User = user;
	return fid, nil;
}

// Connects to a file server and attaches to it as the specified user.
func Mount(net, addr, aname string, user p.User) (*Clnt, *p.Error) {
	clnt, err := Connect(net, addr, 8192+p.IOHdrSz, true);
	if err != nil {
		return nil, err
	}

	fid, err := clnt.Attach(nil, user, aname);
	if err != nil {
		clnt.Unmount();
		return nil, err;
	}

	clnt.Root = fid;
	return clnt, nil;
}

// Closes the connection to the file sever.
func (clnt *Clnt) Unmount() {
	clnt.Lock();
	clnt.err = &p.Error{"connection closed", syscall.ECONNRESET};
	clnt.conn.Close();
	clnt.Unlock();
}