# User Business Module Specification

## Purpose

This domain defines how the existing user example becomes a module-first business module that participates in application assembly through declarations.

## Requirements

### Requirement: User Module Inclusion

The system SHALL include the user business module through the application module list.

#### Scenario: User module is included

- GIVEN the application module list is evaluated by a developer or automated check
- WHEN the included modules are inspected
- THEN the user business module is present in the module list
- AND the user module is identified as a business module

#### Scenario: User module is not included in a custom assembly

- GIVEN a custom application assembly omits the user business module
- WHEN the application is assembled
- THEN the assembly report does not list user as an active module
- AND user HTTP entry points are not exposed by that assembly

### Requirement: User Entry Points Are Declared by the Module

The system SHALL expose user HTTP entry points from the user business module declaration.

#### Scenario: User routes are exposed from module declaration

- GIVEN the user business module is included in application assembly
- WHEN the application is assembled
- THEN the assembly report lists user-owned HTTP entry points
- AND the server route list contains user routes exposed from those entry points

#### Scenario: User entry point adapter is unavailable

- GIVEN the user business module declares HTTP entry points
- AND the application assembly has no platform behavior for HTTP entry points
- WHEN the application is assembled
- THEN assembly fails before the application starts accepting triggers
- AND the failure identifies the user module and the unsupported entry point type

### Requirement: User Storage Capability Is Declared

The system SHALL make the user business module declare a storage capability requirement before startup without naming a concrete storage engine in the business module declaration.

#### Scenario: User storage capability is available

- GIVEN the user business module is included
- AND the application includes an enabled user storage provider that satisfies the user module requirement
- WHEN the application is assembled
- THEN the assembly report marks the user module storage requirement as satisfied
- AND the report identifies the selected storage provider

#### Scenario: User storage capability is missing

- GIVEN the user business module is included
- AND no enabled user storage provider satisfies the user module requirement
- WHEN the application is assembled
- THEN assembly fails before the application starts accepting triggers
- AND the failure identifies the user module and the missing storage capability

### Requirement: User Storage Provider Is Switchable

The system SHALL allow the user module storage provider to be switched without changing the user business module declaration or user HTTP entry points.

#### Scenario: User storage provider is switched

- GIVEN the user business module is included
- AND one enabled user storage provider is selected for the application
- WHEN the selected user storage provider is changed to another enabled user storage provider
- THEN application assembly succeeds with the same user module active
- AND the user HTTP entry points remain unchanged
- AND the assembly report identifies the newly selected storage provider

#### Scenario: User storage provider is not explicitly selected

- GIVEN the user business module is included
- AND no concrete user storage provider is explicitly selected
- WHEN the application is assembled
- THEN application assembly succeeds using a supported default user storage provider
- AND the assembly report identifies the selected storage provider

#### Scenario: User storage provider selection is invalid

- GIVEN the user business module is included
- AND the selected user storage provider is not available to the application
- WHEN the application is assembled
- THEN assembly fails before the application starts accepting triggers
- AND the failure identifies the user module, the storage capability, and the unavailable selected provider

### Requirement: Existing User HTTP Behavior Remains Available

The system SHALL keep the existing user HTTP behavior available after the user module migration.

#### Scenario: Existing user route is reachable

- GIVEN the application is assembled with the user business module included
- WHEN an existing user route is invoked through the server
- THEN the request reaches user behavior through the module-first path
- AND the response uses the existing success or error envelope format

#### Scenario: Existing user route receives invalid input

- GIVEN the application is assembled with the user business module included
- WHEN an existing user route receives invalid input
- THEN the request is rejected through the existing error envelope format
- AND the response does not bypass platform error handling

### Requirement: OpenAPI And Generated Output Do Not Drive User Assembly

The system SHALL assemble the user business module without requiring OpenAPI or generated application wiring.

#### Scenario: User module assembles without OpenAPI input

- GIVEN no OpenAPI input is provided to application assembly
- AND the user business module is included
- WHEN the application is assembled
- THEN user module assembly succeeds using module declarations
- AND user HTTP entry points are available

#### Scenario: Legacy generated files exist

- GIVEN legacy generated or layered user files exist in the project
- WHEN the application is assembled
- THEN those files do not make user active unless the user business module is explicitly included
- AND generated output does not override the user module's declared capabilities or entry points
