// Copyright (C) 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import "fmt"

// cmdSet is the namespace for a command identifier.
type cmdSet uint8

// cmdID is a command in a command set.
type cmdID uint8

type Cmd struct {
	set cmdSet
	id  cmdID
}

func (c Cmd) String() string {
	return fmt.Sprintf("%v.%v", c.set, cmdNames[c])
}

const (
	cmdSetVirtualMachine       = cmdSet(1)
	cmdSetReferenceType        = cmdSet(2)
	cmdSetClassType            = cmdSet(3)
	cmdSetArrayType            = cmdSet(4)
	cmdSetInterfaceType        = cmdSet(5)
	cmdSetMethod               = cmdSet(6)
	cmdSetField                = cmdSet(8)
	cmdSetObjectReference      = cmdSet(9)
	cmdSetStringReference      = cmdSet(10)
	cmdSetThreadReference      = cmdSet(11)
	cmdSetThreadGroupReference = cmdSet(12)
	cmdSetArrayReference       = cmdSet(13)
	cmdSetClassLoaderReference = cmdSet(14)
	cmdSetEventRequest         = cmdSet(15)
	cmdSetStackFrame           = cmdSet(16)
	cmdSetClassObjectReference = cmdSet(17)
	cmdSetEvent                = cmdSet(64)
)

var (
	CmdVirtualMachineVersion               = Cmd{cmdSetVirtualMachine, 1}
	CmdVirtualMachineClassesBySignature    = Cmd{cmdSetVirtualMachine, 2}
	CmdVirtualMachineAllClasses            = Cmd{cmdSetVirtualMachine, 3}
	CmdVirtualMachineAllThreads            = Cmd{cmdSetVirtualMachine, 4}
	CmdVirtualMachineTopLevelThreadGroups  = Cmd{cmdSetVirtualMachine, 5}
	CmdVirtualMachineDispose               = Cmd{cmdSetVirtualMachine, 6}
	CmdVirtualMachineIDSizes               = Cmd{cmdSetVirtualMachine, 7}
	CmdVirtualMachineSuspend               = Cmd{cmdSetVirtualMachine, 8}
	CmdVirtualMachineResume                = Cmd{cmdSetVirtualMachine, 9}
	CmdVirtualMachineExit                  = Cmd{cmdSetVirtualMachine, 10}
	CmdVirtualMachineCreateString          = Cmd{cmdSetVirtualMachine, 11}
	CmdVirtualMachineCapabilities          = Cmd{cmdSetVirtualMachine, 12}
	CmdVirtualMachineClassPaths            = Cmd{cmdSetVirtualMachine, 13}
	CmdVirtualMachineDisposeObjects        = Cmd{cmdSetVirtualMachine, 14}
	CmdVirtualMachineHoldEvents            = Cmd{cmdSetVirtualMachine, 15}
	CmdVirtualMachineReleaseEvents         = Cmd{cmdSetVirtualMachine, 16}
	CmdVirtualMachineCapabilitiesNew       = Cmd{cmdSetVirtualMachine, 17}
	CmdVirtualMachineRedefineClasses       = Cmd{cmdSetVirtualMachine, 18}
	CmdVirtualMachineSetDefaultStratum     = Cmd{cmdSetVirtualMachine, 19}
	CmdVirtualMachineAllClassesWithGeneric = Cmd{cmdSetVirtualMachine, 20}
	CmdVirtualMachineInstanceCounts        = Cmd{cmdSetVirtualMachine, 21}

	CmdReferenceTypeSignature            = Cmd{cmdSetReferenceType, 1}
	CmdReferenceTypeClassLoader          = Cmd{cmdSetReferenceType, 2}
	CmdReferenceTypeModifiers            = Cmd{cmdSetReferenceType, 3}
	CmdReferenceTypeFields               = Cmd{cmdSetReferenceType, 4}
	CmdReferenceTypeMethods              = Cmd{cmdSetReferenceType, 5}
	CmdReferenceTypeGetValues            = Cmd{cmdSetReferenceType, 6}
	CmdReferenceTypeSourceFile           = Cmd{cmdSetReferenceType, 7}
	CmdReferenceTypeNestedTypes          = Cmd{cmdSetReferenceType, 8}
	CmdReferenceTypeStatus               = Cmd{cmdSetReferenceType, 9}
	CmdReferenceTypeInterfaces           = Cmd{cmdSetReferenceType, 10}
	CmdReferenceTypeClassObject          = Cmd{cmdSetReferenceType, 11}
	CmdReferenceTypeSourceDebugExtension = Cmd{cmdSetReferenceType, 12}
	CmdReferenceTypeSignatureWithGeneric = Cmd{cmdSetReferenceType, 13}
	CmdReferenceTypeFieldsWithGeneric    = Cmd{cmdSetReferenceType, 14}
	CmdReferenceTypeMethodsWithGeneric   = Cmd{cmdSetReferenceType, 15}
	CmdReferenceTypeInstances            = Cmd{cmdSetReferenceType, 16}
	CmdReferenceTypeClassFileVersion     = Cmd{cmdSetReferenceType, 17}
	CmdReferenceTypeConstantPool         = Cmd{cmdSetReferenceType, 18}

	CmdClassTypeSuperclass   = Cmd{cmdSetClassType, 1}
	CmdClassTypeSetValues    = Cmd{cmdSetClassType, 2}
	CmdClassTypeInvokeMethod = Cmd{cmdSetClassType, 3}
	CmdClassTypeNewInstance  = Cmd{cmdSetClassType, 4}

	CmdArrayTypeNewInstance = Cmd{cmdSetArrayType, 1}

	CmdInterfaceTypeInvokeMethod = Cmd{cmdSetInterfaceType, 1}

	CmdMethodTypeLineTable                = Cmd{cmdSetMethod, 1}
	CmdMethodTypeVariableTable            = Cmd{cmdSetMethod, 2}
	CmdMethodTypeBytecodes                = Cmd{cmdSetMethod, 3}
	CmdMethodTypeIsObsolete               = Cmd{cmdSetMethod, 4}
	CmdMethodTypeVariableTableWithGeneric = Cmd{cmdSetMethod, 5}

	CmdObjectReferenceReferenceType     = Cmd{cmdSetObjectReference, 1}
	CmdObjectReferenceGetValues         = Cmd{cmdSetObjectReference, 2}
	CmdObjectReferenceSetValues         = Cmd{cmdSetObjectReference, 3}
	CmdObjectReferenceMonitorInfo       = Cmd{cmdSetObjectReference, 5}
	CmdObjectReferenceInvokeMethod      = Cmd{cmdSetObjectReference, 6}
	CmdObjectReferenceDisableCollection = Cmd{cmdSetObjectReference, 7}
	CmdObjectReferenceEnableCollection  = Cmd{cmdSetObjectReference, 8}
	CmdObjectReferenceIsCollected       = Cmd{cmdSetObjectReference, 9}
	CmdObjectReferringObjects           = Cmd{cmdSetObjectReference, 10}

	CmdStringReferenceValue = Cmd{cmdSetStringReference, 1}

	CmdThreadReferenceName                    = Cmd{cmdSetThreadReference, 1}
	CmdThreadReferenceSuspend                 = Cmd{cmdSetThreadReference, 2}
	CmdThreadReferenceResume                  = Cmd{cmdSetThreadReference, 3}
	CmdThreadReferenceStatus                  = Cmd{cmdSetThreadReference, 4}
	CmdThreadReferenceThreadGroup             = Cmd{cmdSetThreadReference, 5}
	CmdThreadReferenceFrames                  = Cmd{cmdSetThreadReference, 6}
	CmdThreadReferenceFrameCount              = Cmd{cmdSetThreadReference, 7}
	CmdThreadReferenceOwnedMonitors           = Cmd{cmdSetThreadReference, 8}
	CmdThreadReferenceCurrentContendedMonitor = Cmd{cmdSetThreadReference, 9}
	CmdThreadReferenceStop                    = Cmd{cmdSetThreadReference, 10}
	CmdThreadReferenceInterrupt               = Cmd{cmdSetThreadReference, 11}
	CmdThreadReferenceSuspendCount            = Cmd{cmdSetThreadReference, 12}

	CmdThreadGroupReferenceName     = Cmd{cmdSetThreadGroupReference, 1}
	CmdThreadGroupReferenceParent   = Cmd{cmdSetThreadGroupReference, 2}
	CmdThreadGroupReferenceChildren = Cmd{cmdSetThreadGroupReference, 3}

	CmdArrayReferenceLength    = Cmd{cmdSetArrayReference, 1}
	CmdArrayReferenceGetValues = Cmd{cmdSetArrayReference, 2}
	CmdArrayReferenceSetValues = Cmd{cmdSetArrayReference, 3}

	CmdClassLoaderReferenceVisibleClasses = Cmd{cmdSetClassLoaderReference, 1}

	CmdEventRequestSet                 = Cmd{cmdSetEventRequest, 1}
	CmdEventRequestClear               = Cmd{cmdSetEventRequest, 2}
	CmdEventRequestClearAllBreakpoints = Cmd{cmdSetEventRequest, 3}

	CmdStackFrameGetValues  = Cmd{cmdSetStackFrame, 1}
	CmdStackFrameSetValues  = Cmd{cmdSetStackFrame, 2}
	CmdStackFrameThisObject = Cmd{cmdSetStackFrame, 3}
	CmdStackFramePopFrames  = Cmd{cmdSetStackFrame, 4}

	CmdClassObjectReferenceReflectedType = Cmd{cmdSetClassObjectReference, 1}

	CmdEventComposite = Cmd{cmdSetEvent, 1}
)

var cmdNames = map[Cmd]string{}

func init() {
	register := func(c Cmd, n string) {
		if _, e := cmdNames[c]; e {
			panic("command already registered")
		}
		cmdNames[c] = n
	}
	register(CmdVirtualMachineVersion, "Version")
	register(CmdVirtualMachineClassesBySignature, "ClassesBySignature")
	register(CmdVirtualMachineAllClasses, "AllClasses")
	register(CmdVirtualMachineAllThreads, "AllThreads")
	register(CmdVirtualMachineTopLevelThreadGroups, "TopLevelThreadGroups")
	register(CmdVirtualMachineDispose, "Dispose")
	register(CmdVirtualMachineIDSizes, "IDSizes")
	register(CmdVirtualMachineSuspend, "Suspend")
	register(CmdVirtualMachineResume, "Resume")
	register(CmdVirtualMachineExit, "Exit")
	register(CmdVirtualMachineCreateString, "CreateString")
	register(CmdVirtualMachineCapabilities, "Capabilities")
	register(CmdVirtualMachineClassPaths, "ClassPaths")
	register(CmdVirtualMachineDisposeObjects, "DisposeObjects")
	register(CmdVirtualMachineHoldEvents, "HoldEvents")
	register(CmdVirtualMachineReleaseEvents, "ReleaseEvents")
	register(CmdVirtualMachineCapabilitiesNew, "CapabilitiesNew")
	register(CmdVirtualMachineRedefineClasses, "RedefineClasses")
	register(CmdVirtualMachineSetDefaultStratum, "SetDefaultStratum")
	register(CmdVirtualMachineAllClassesWithGeneric, "AllClassesWithGeneric")
	register(CmdVirtualMachineInstanceCounts, "InstanceCounts")
	register(CmdReferenceTypeClassFileVersion, "ClassFileVersion")

	register(CmdReferenceTypeSignature, "Signature")
	register(CmdReferenceTypeClassLoader, "ClassLoader")
	register(CmdReferenceTypeModifiers, "Modifiers")
	register(CmdReferenceTypeFields, "Fields")
	register(CmdReferenceTypeMethods, "Methods")
	register(CmdReferenceTypeGetValues, "GetValues")
	register(CmdReferenceTypeSourceFile, "SourceFile")
	register(CmdReferenceTypeNestedTypes, "NestedTypes")
	register(CmdReferenceTypeStatus, "Status")
	register(CmdReferenceTypeInterfaces, "Interfaces")
	register(CmdReferenceTypeClassObject, "ClassObject")
	register(CmdReferenceTypeSourceDebugExtension, "SourceDebugExtension")
	register(CmdReferenceTypeSignatureWithGeneric, "SignatureWithGeneric")
	register(CmdReferenceTypeFieldsWithGeneric, "FieldsWithGeneric")
	register(CmdReferenceTypeMethodsWithGeneric, "MethodsWithGeneric")
	register(CmdReferenceTypeInstances, "Instances")
	register(CmdReferenceTypeConstantPool, "ConstantPool")

	register(CmdClassTypeSuperclass, "Superclass")
	register(CmdClassTypeSetValues, "SetValues")
	register(CmdClassTypeInvokeMethod, "InvokeMethod")
	register(CmdClassTypeNewInstance, "NewInstance")

	register(CmdArrayTypeNewInstance, "NewInstance")

	register(CmdMethodTypeLineTable, "LineTable")
	register(CmdMethodTypeVariableTable, "VariableTable")
	register(CmdMethodTypeBytecodes, "Bytecodes")
	register(CmdMethodTypeIsObsolete, "IsObsolete")
	register(CmdMethodTypeVariableTableWithGeneric, "VariableTableWithGeneric")

	register(CmdObjectReferenceReferenceType, "ReferenceType")
	register(CmdObjectReferenceGetValues, "GetValues")
	register(CmdObjectReferenceSetValues, "SetValues")
	register(CmdObjectReferenceMonitorInfo, "MonitorInfo")
	register(CmdObjectReferenceInvokeMethod, "InvokeMethod")
	register(CmdObjectReferenceDisableCollection, "DisableCollection")
	register(CmdObjectReferenceEnableCollection, "EnableCollection")
	register(CmdObjectReferenceIsCollected, "IsCollected")
	register(CmdObjectReferringObjects, "ReferringObjects")

	register(CmdStringReferenceValue, "Value")

	register(CmdThreadReferenceName, "Name")
	register(CmdThreadReferenceSuspend, "Suspend")
	register(CmdThreadReferenceResume, "Resume")
	register(CmdThreadReferenceStatus, "Status")
	register(CmdThreadReferenceThreadGroup, "ThreadGroup")
	register(CmdThreadReferenceFrames, "Frames")
	register(CmdThreadReferenceFrameCount, "FrameCount")
	register(CmdThreadReferenceOwnedMonitors, "OwnedMonitors")
	register(CmdThreadReferenceCurrentContendedMonitor, "CurrentContendedMonitor")
	register(CmdThreadReferenceStop, "Stop")
	register(CmdThreadReferenceInterrupt, "Interrupt")
	register(CmdThreadReferenceSuspendCount, "SuspendCount")

	register(CmdThreadGroupReferenceName, "Name")
	register(CmdThreadGroupReferenceParent, "Parent")
	register(CmdThreadGroupReferenceChildren, "Children")

	register(CmdArrayReferenceLength, "Length")
	register(CmdArrayReferenceGetValues, "GetValues")
	register(CmdArrayReferenceSetValues, "SetValues")

	register(CmdClassLoaderReferenceVisibleClasses, "VisibleClasses")

	register(CmdEventRequestSet, "Set")
	register(CmdEventRequestClear, "Clear")
	register(CmdEventRequestClearAllBreakpoints, "ClearAllBreakpoints")

	register(CmdStackFrameGetValues, "GetValues")
	register(CmdStackFrameSetValues, "SetValues")
	register(CmdStackFrameThisObject, "ThisObject")
	register(CmdStackFramePopFrames, "PopFrames")

	register(CmdClassObjectReferenceReflectedType, "ReflectedType")

	register(CmdEventComposite, "Composite")
}
