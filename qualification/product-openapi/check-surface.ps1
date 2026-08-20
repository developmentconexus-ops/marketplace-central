[CmdletBinding()]
param(
    [string]$OadPath = (Join-Path $PSScriptRoot '../../contracts/api/product/openapi.yaml')
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

function Fail([string]$Message) {
    Write-Error $Message
    exit 1
}

$expectedRows = @(
    'GetCurrentAccessContext|Q|authenticated|H,A,S',
    'ListOrganizationMembers|Q|access.read|H',
    'ListAccessRoles|Q|access.read|H',
    'AssignAccessRole|C|access.manage|H',
    'RevokeAccessRole|C|access.manage|H',
    'ListMarketplaceInstallations|Q|portfolio.read|H,A,S',
    'GetMarketplaceInstallation|Q|portfolio.read|H,A,S',
    'CreateMarketplaceInstallation|C|portfolio.manage|H',
    'UpdateMarketplaceInstallationConfiguration|C|portfolio.manage|H',
    'DeactivateMarketplaceInstallation|C|portfolio.manage|H',
    'ListSellingEntities|Q|portfolio.read|H,A,S',
    'SearchSourceProductsForMarketplace|Q|readiness.read|H,A,S',
    'GetProductChannelReadiness|Q|readiness.read|H,A,S',
    'GetPublicationRequirements|Q|readiness.read|H,A,S',
    'ResolveProductChannelCorrespondence|C|readiness.manage|H,A',
    'ClearProductChannelCorrespondence|C|readiness.manage|H,A',
    'ListMarketplaceListings|Q|offering.read|H,A,S',
    'GetMarketplaceListing|Q|offering.read|H,A,S',
    'ListListingIntents|Q|offering.read|H,A,S',
    'GetListingIntent|Q|offering.read|H,A,S',
    'CreateListingIntentDraft|C|listing.manage|H,A',
    'UpdateListingIntentDraft|C|listing.manage|H,A',
    'DiscardListingIntentDraft|C|listing.manage|H,A',
    'SubmitListingIntent|C|listing.manage|H,A',
    'CreateListingIntentMedia|C|listing.manage|H,A',
    'ListPriceIntents|Q|offering.read|H,A,S',
    'GetPriceIntent|Q|offering.read|H,A,S',
    'CreatePriceIntent|C|price.manage|H,A',
    'ListSellableAvailability|Q|availability.read|H,A,S',
    'GetSellableAvailability|Q|availability.read|H,A,S',
    'ListInventorySources|Q|availability.read|H,A,S',
    'GetInventorySource|Q|availability.read|H,A,S',
    'CreateInventorySource|C|availability.manage|H',
    'UpdateInventorySource|C|availability.manage|H',
    'DeactivateInventorySource|C|availability.manage|H',
    'GetEffectiveAvailabilityAllocationScopePolicy|Q|availability.read|H,A,S',
    'UpdateAvailabilityAllocationScopePolicy|C|availability.manage|H',
    'ListCompetitivePositions|Q|market.read|H,A,S',
    'GetCompetitivePosition|Q|market.read|H,A,S',
    'ListComparableOffers|Q|market.read|H,A,S',
    'ListExpectedEconomics|Q|economics.read|H,A,S',
    'GetExpectedEconomics|Q|economics.read|H,A,S',
    'EvaluatePriceScenario|C|economics.read|H,A,S',
    'ListSaleEconomics|Q|economics.read|H,A,S',
    'GetSaleEconomics|Q|economics.read|H,A,S',
    'GetEconomicPerformanceSummary|Q|economics.read|H,A,S',
    'GetCommercialPolicy|Q|economics.read|H,A,S',
    'UpdateCommercialPolicy|C|economics.policy.manage|H',
    'ListEconomicAttributions|Q|economics.read|H,A,S',
    'GetEconomicAttribution|Q|economics.read|H,A,S',
    'ResolveEconomicAttribution|C|economics.reconcile|H',
    'ListAuthorizationDecisions|Q|governance.read|H,A,S',
    'GetAuthorizationDecision|Q|governance.read|H,A,S',
    'CreateAuthorizationDecision|C|governance.decide|H',
    'ListAuthorizationDelegations|Q|governance.manage|H',
    'EstablishAuthorizationDelegation|C|governance.manage|H',
    'UpdateAuthorizationDelegation|C|governance.manage|H',
    'RevokeAuthorizationDelegation|C|governance.manage|H',
    'ListMarketplaceSales|Q|sales.read|H,A,S',
    'GetMarketplaceSale|Q|sales.read|H,A,S',
    'ResolveSaleSellingEntityAttribution|C|sales.manage|H',
    'ListBusinessOrderIntents|Q|materialization.read|H,A,S',
    'GetBusinessOrderIntent|Q|materialization.read|H,A,S',
    'GetBusinessSystemPartyResolution|Q|materialization.read|H,A,S',
    'ResolveBusinessSystemPartyResolution|C|materialization.resolve|H',
    'GetDestinationRealization|Q|materialization.read|H,A,S',
    'ListInvoicingIntents|Q|materialization.read|H,A,S',
    'GetInvoicingIntent|Q|materialization.read|H,A,S',
    'ListFulfillmentExecutions|Q|fulfillment.read|H,A,S',
    'GetFulfillmentExecution|Q|fulfillment.read|H,A,S',
    'ListFulfillmentNodes|Q|fulfillment.read|H,A,S',
    'GetFulfillmentNode|Q|fulfillment.read|H,A,S',
    'CreateFulfillmentNode|C|fulfillment.manage|H',
    'UpdateFulfillmentNode|C|fulfillment.manage|H',
    'DeactivateFulfillmentNode|C|fulfillment.manage|H',
    'GetFulfillmentOperatingTargets|Q|fulfillment.read|H,A,S',
    'UpdateFulfillmentOperatingTargets|C|fulfillment.manage|H',
    'RecordSeparation|C|fulfillment.execute|H',
    'RecordPhysicalConference|C|fulfillment.execute|H,S',
    'RecordPacking|C|fulfillment.execute|H',
    'RecordDispatchHandoff|C|fulfillment.execute|H,S',
    'ListFulfillmentArtifacts|Q|fulfillment.execute|H,A,S',
    'GetFulfillmentArtifact|Q|fulfillment.execute|H,A,S',
    'ListShipments|Q|fulfillment.read|H,A,S',
    'GetShipment|Q|fulfillment.read|H,A,S',
    'ListPostSaleResolutions|Q|post_sale.read|H,A,S',
    'GetPostSaleResolution|Q|post_sale.read|H,A,S',
    'CreatePostSaleResolution|C|post_sale.manage|H,A',
    'ListWork|Q|work.read|H,A,S',
    'GetWork|Q|work.read|H,A,S',
    'AssignWork|C|work.manage|H,A',
    'ClearWorkAssignment|C|work.manage|H,A',
    'HoldWork|C|work.manage|H,A',
    'ResumeWork|C|work.manage|H,A',
    'EscalateWork|C|work.manage|H,A'
)

$expectedPermissions = @(
    'access.read','access.manage','portfolio.read','portfolio.manage','readiness.read','readiness.manage',
    'offering.read','listing.manage','price.manage','availability.read','availability.manage','market.read',
    'economics.read','economics.policy.manage','economics.reconcile','governance.read','governance.decide',
    'governance.manage','sales.read','sales.manage','materialization.read','materialization.resolve',
    'fulfillment.read','fulfillment.execute','fulfillment.manage','post_sale.read','post_sale.manage',
    'work.read','work.manage'
)

function Get-Surface([string]$Text) {
    $pattern = 'operationId:\s*([A-Z][A-Za-z0-9]*)\s*,\s*x-mpc-operation-class:\s*(Q|C)\s*,\s*x-mpc-required-permission:\s*([^,\s}]+)\s*,\s*x-mpc-principal-kinds:\s*\[([^\]]+)\]'
    $matches = [regex]::Matches($Text, $pattern)
    $rows = @()
    foreach ($match in $matches) {
        $principals = (($match.Groups[4].Value -split ',') | ForEach-Object { $_.Trim() }) -join ','
        $rows += "$($match.Groups[1].Value)|$($match.Groups[2].Value)|$($match.Groups[3].Value)|$principals"
    }
    return @($rows)
}

function Test-Surface([string]$Text, [switch]$Quiet) {
    $errors = @()
    if (-not $Text.Contains('openapi: 3.1.2')) { $errors += 'OpenAPI feature level is not exactly 3.1.2.' }
    if (-not $Text.Contains('jsonSchemaDialect: https://spec.openapis.org/oas/3.1/dialect/base')) { $errors += 'JSON Schema dialect is not the OAS 3.1 base dialect.' }
    if ($Text -match '(?m)^servers\s*:') { $errors += 'Environment/runtime servers must be omitted from the source OAD.' }
    if (-not $Text.Contains('MpcBearerAuth:')) { $errors += 'MpcBearerAuth scheme is missing.' }

    $rows = @(Get-Surface $Text)
    if ($rows.Count -ne 95) { $errors += "Expected 95 operation rows; found $($rows.Count)." }
    if (@($rows | Sort-Object -Unique).Count -ne $rows.Count) { $errors += 'Duplicate operation enforcement rows detected.' }

    $missing = @($expectedRows | Where-Object { $_ -notin $rows })
    $extra = @($rows | Where-Object { $_ -notin $expectedRows })
    if ($missing.Count -gt 0) { $errors += ('Missing/mismatched W4 rows: ' + ($missing -join '; ')) }
    if ($extra.Count -gt 0) { $errors += ('Unexpected W4 rows: ' + ($extra -join '; ')) }

    $permissionMatches = [regex]::Matches($Text, 'x-mpc-required-permission:\s*([^,\s}]+)')
    $permissions = @($permissionMatches | ForEach-Object { $_.Groups[1].Value } | Where-Object { $_ -ne 'authenticated' } | Sort-Object -Unique)
    if ($permissions.Count -ne 29) { $errors += "Expected 29 stored ordinary Permissions; found $($permissions.Count)." }
    if (@($expectedPermissions | Where-Object { $_ -notin $permissions }).Count -gt 0 -or @($permissions | Where-Object { $_ -notin $expectedPermissions }).Count -gt 0) {
        $errors += 'Ordinary Permission vocabulary differs from canonical W4.'
    }

    $principalTokens = @([regex]::Matches($Text, 'x-mpc-principal-kinds:\s*\[([^\]]+)\]') | ForEach-Object {
        $_.Groups[1].Value -split ',' | ForEach-Object { $_.Trim() }
    } | Sort-Object -Unique)
    if (($principalTokens -join ',') -ne 'A,H,S') { $errors += ('Principal kind vocabulary must be exactly H/A/S; found ' + ($principalTokens -join ',') + '.') }

    $operationIds = @([regex]::Matches($Text, 'operationId:\s*([A-Z][A-Za-z0-9]*)') | ForEach-Object { $_.Groups[1].Value })
    if (@($operationIds | Sort-Object -Unique).Count -ne 95) { $errors += 'operationId values are not exactly 95 unique PascalCase identifiers.' }
    $listSearchCount = @($operationIds | Where-Object { $_ -match '^(List|Search)' }).Count
    if ($listSearchCount -ne 26) { $errors += "Expected 26 List/Search operations; found $listSearchCount." }

    $pathLines = @($Text -split "`n" | Where-Object { $_ -match '^  /' })
    foreach ($line in $pathLines) {
        $path = $line.Trim()
        $path = $path.Substring(0, $path.Length - 1)
        if ($path -ne '/access-context' -and -not $path.StartsWith('/organizations/{organization_id}/', [StringComparison]::Ordinal)) {
            $errors += "Non-canonical Product root: $path"
        }
    }

    foreach ($forbidden in @('/v1','/providers','/integrations','/webhooks','/oauth','/mutations','ngrok')) {
        if ($Text.Contains($forbidden, [StringComparison]::OrdinalIgnoreCase)) { $errors += "Forbidden Product OAD token present: $forbidden" }
    }

    $physical = @([regex]::Matches($Text, 'x-mpc-physical-qualification:\s*([^,\s}]+)') | ForEach-Object { $_.Groups[1].Value } | Sort-Object)
    if (($physical -join ',') -ne 'fulfillment.dispatch-handoff,fulfillment.physical-conference') { $errors += 'Physical qualification markers differ from the two canonical checkpoints.' }

    if (-not $Quiet -and $errors.Count -gt 0) { $errors | ForEach-Object { Write-Error $_ } }
    return $errors.Count -eq 0
}

if (-not (Test-Path -LiteralPath $OadPath -PathType Leaf)) { Fail "Product OAD missing: $OadPath" }
$oad = Get-Content -LiteralPath $OadPath -Raw
if (-not (Test-Surface $oad)) { exit 1 }

$negativeFixture = $oad.Replace('operationId: GetCurrentAccessContext', 'operationId: GetCurrentAccessContextBROKEN')
if (Test-Surface $negativeFixture -Quiet) { Fail 'Negative control failed: a corrupted operationId was accepted.' }

Write-Host 'product_oad_surface_operations: 95/95'
Write-Host 'product_oad_surface_permissions: 29/29 + authenticated special condition'
Write-Host 'product_oad_surface_principal_kinds: H/A/S only'
Write-Host 'product_oad_surface_list_search: 26/26'
Write-Host 'product_oad_surface_negative_control: PASS'
Write-Host 'product_oad_surface: PASS'
