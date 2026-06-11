# Module-First Architecture Specification

## Purpose

This domain defines how developers add business behavior through business modules while platform and capability modules provide reusable runtime behavior.

## Requirements

### Requirement: Business Modules Are the Source of Truth

The system SHALL allow a developer to add or change business behavior by working primarily within a business module.

#### Scenario: Developer adds a complete business module

- GIVEN a developer has a business module with declared required capabilities and declared entry points
- WHEN the application is assembled with that business module included
- THEN the module's declared entry points are available through the application
- AND the developer does not need to edit unrelated business modules to expose those entry points

#### Scenario: Business module is incomplete

- GIVEN a developer has a business module with no declared entry points
- WHEN the application is assembled with that business module included
- THEN the application reports that the module contributes no entry points
- AND the report identifies the affected business module

### Requirement: Platform Owns Shared Runtime Behavior

The system SHALL provide shared runtime behavior through platform modules rather than requiring each business module to repeat startup, configuration, logging, error handling, response handling, lifecycle, health check, or observability behavior.

#### Scenario: Shared runtime behavior is available to a business module

- GIVEN a business module is included in an application with platform modules enabled
- WHEN the business module handles an incoming entry point trigger
- THEN shared runtime behavior is applied consistently around that trigger
- AND the business module is not required to define its own startup, configuration loading, logging setup, error translation, response envelope, lifecycle runner, health aggregation, or observability setup

#### Scenario: Required platform behavior is unavailable

- GIVEN a business module is included in an application that cannot provide a required shared runtime behavior
- WHEN the application is assembled
- THEN assembly fails before the application starts accepting entry point triggers
- AND the failure identifies the missing platform behavior and the affected business module

### Requirement: Capabilities Are Declared by Business Modules

The system SHALL allow each business module to declare the reusable capabilities it needs before the application starts.

#### Scenario: Required capability is enabled

- GIVEN a business module declares a required capability
- AND the application includes an enabled capability module that satisfies that requirement
- WHEN the application is assembled
- THEN the business module is connected to that capability
- AND the application reports that the requirement is satisfied

#### Scenario: Required capability is missing

- GIVEN a business module declares a required capability
- AND the application does not include an enabled capability module that satisfies that requirement
- WHEN the application is assembled
- THEN assembly fails before the application starts accepting entry point triggers
- AND the failure identifies the missing capability and the affected business module

### Requirement: Capability Modules Are Optional and Observable

The system SHALL make optional infrastructure capabilities observable as enabled, disabled, healthy, or unavailable without requiring business modules to inspect infrastructure setup directly.

#### Scenario: Optional capability is disabled

- GIVEN an application configuration disables an optional capability
- AND no included business module requires that capability
- WHEN the application is assembled
- THEN the application starts without that capability
- AND the capability status is observable as disabled

#### Scenario: Enabled capability is unavailable

- GIVEN an application configuration enables a capability
- AND the capability cannot become available during startup
- WHEN the application is assembled
- THEN assembly fails before the application starts accepting entry point triggers
- AND the failure identifies the unavailable capability

### Requirement: Entry Points Are Declared by Business Modules

The system SHALL allow business modules to declare the external triggers that reach their business behavior.

#### Scenario: Entry point declaration is accepted

- GIVEN a business module declares an entry point with a trigger type and business action
- WHEN the application is assembled
- THEN the entry point is exposed through the matching platform behavior
- AND the application can list the entry point as belonging to that business module

#### Scenario: Entry point declaration cannot be exposed

- GIVEN a business module declares an entry point
- AND the application lacks the platform behavior needed to expose that entry point type
- WHEN the application is assembled
- THEN assembly fails before the application starts accepting entry point triggers
- AND the failure identifies the unsupported entry point type and the affected business module

### Requirement: Normal Business Work Avoids Central Wiring Changes

The system SHALL support normal business feature work without requiring developers to edit central application wiring for each new business module, capability use, or entry point.

#### Scenario: Business module is added through module inclusion

- GIVEN a developer has created a business module that declares its required capabilities and entry points
- WHEN the developer includes the business module in the application module list
- THEN the platform composes the module's capabilities and entry points during assembly
- AND no additional central wiring changes are required for the declared capabilities and entry points

#### Scenario: Module inclusion is omitted

- GIVEN a developer has created a business module
- AND the business module is not included in the application module list
- WHEN the application is assembled
- THEN the module's entry points are not exposed
- AND the application module list does not report the module as active

### Requirement: OpenAPI Is Not the Architecture Source of Truth

The system SHALL not require an OpenAPI document to define business modules, capabilities, or application composition.

#### Scenario: Application assembles without OpenAPI input

- GIVEN business modules, platform modules, and capability modules are available
- AND no OpenAPI document is provided
- WHEN the application is assembled
- THEN application assembly can complete using module declarations
- AND business module entry points can be exposed without OpenAPI input

#### Scenario: OpenAPI input is present

- GIVEN an OpenAPI document is present in the project
- WHEN the application architecture is evaluated
- THEN the OpenAPI document is treated as optional supporting input
- AND it does not override business module declarations, capability declarations, or application composition decisions

### Requirement: Code Generation Is Optional Assistance

The system SHALL not require code generation as the primary workflow for creating or changing business behavior.

#### Scenario: Developer creates business behavior manually

- GIVEN a developer creates a business module, declares required capabilities, and declares entry points without running a generator
- WHEN the application is assembled with that module included
- THEN the module can participate in application composition
- AND the module can expose its declared entry points

#### Scenario: Generated output exists

- GIVEN generated output exists in the project
- WHEN the application architecture is evaluated
- THEN generated output is treated as editable project input only when included through module declarations
- AND generated output does not become the source of truth for business behavior, capability requirements, or application composition
