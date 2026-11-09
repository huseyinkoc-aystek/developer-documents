# P2-E14: e-Fatura Integration — Production Setup Guide

**Date:** February 2, 2026  
**Status:** Ready for Production Deployment  
**Target:** GIB (Turkish Revenue Administration) e-Fatura System

---

## 📋 Pre-Launch Checklist

### 1. GIB Registration & Credentials

**Step 1: Register with GIB**
- [ ] Go to https://efatura.gib.gov.tr
- [ ] Complete company registration with tax ID (VKN)
- [ ] Request e-invoice authorization from GIB
- [ ] Wait for approval (typically 2-3 working days)

**Step 2: Obtain Production Credentials**
- [ ] Download GIB production certificates (SSL/TLS)
- [ ] Obtain username for e-invoice service
- [ ] Obtain password for e-invoice service
- [ ] Note: Production credentials are different from test/sandbox

**Step 3: Digital Signature Setup**
- [ ] Obtain e-signature certificate (from authorized provider)
- [ ] Download certificate file (.cer, .pfx)
- [ ] Download private key file (.key)
- [ ] Store securely (never commit to repo)

### 2. Environment Variable Configuration

**Production Environment Variables:**

```bash
# GIB Configuration (Production)
export GIB_ENVIRONMENT=prod
export GIB_BASE_URL=https://efatura.gib.gov.tr/v1/soap
export GIB_USERNAME=your_production_username
export GIB_PASSWORD=your_production_password
export GIB_VKN_TCKN=12345678900              # Your company's tax ID
export GIB_ALIAS=urn:mail:invoices@company.com.tr  # GIB registered alias

# Certificate Paths
export GIB_CERT_PATH=/secure/path/to/certificate.cer
export GIB_PRIVATE_KEY_PATH=/secure/path/to/private.key

# Logging (optional)
export GIB_LOG_LEVEL=info  # or debug for troubleshooting
```

**Important:**
- Store certificates in secure location outside repo (e.g., Kubernetes secrets, AWS Secrets Manager)
- Use strong permissions (chmod 600) on private key files
- Rotate credentials annually

### 3. Sandbox Validation (Before Production)

**Step 1: Test with GIB Sandbox**

```bash
# Set sandbox credentials
export GIB_ENVIRONMENT=test
export GIB_BASE_URL=https://efaturatest.izibiz.com.tr
export GIB_USERNAME=sandbox_username
export GIB_PASSWORD=sandbox_password
export GIB_VKN_TCKN=1234567890  # Test VKN
```

**Step 2: Run Integration Tests**

```bash
# Run sandbox test suite
cd /workspace/go/aystek-core
go test -v ./internal/modules/billing/invoices/adapters -run TestGIBSandbox

# Expected output:
# ✅ TestGIBSandbox_Integration - Sandbox connectivity
# ✅ TestGIBSandbox_SendInvoice - Invoice submission
# ✅ TestGIBSandbox_StatusMapping - Status normalization
# ✅ TestGIBSandbox_InvalidInvoice - Error handling
# ✅ TestGIBSandbox_NetworkResilience - Network failure recovery
# ✅ TestGIBSandbox_CancelOperation - Invoice cancellation
# ✅ TestGIBSandbox_QueryStatus - Status querying
# ✅ TestGIBSandbox_EnvironmentConfiguration - Config loading
```

**Step 3: Validate Sandbox Operations**

```bash
# Send test invoice
curl -X POST http://localhost:8080/admin/invoices \
  -H "Authorization: Bearer token" \
  -H "Content-Type: application/json" \
  -d '{
    "customer_name": "Test Customer",
    "lines": [{
      "name": "Test Item",
      "quantity": 1,
      "unit_price": 100.0
    }],
    "currency": "TRY"
  }'

# Response should include invoice ID with eissue metadata
# Example: "id": "507f1f77bcf86cd799439011"

# Query e-invoice status
curl -X GET http://localhost:8080/admin/invoices/507f1f77bcf86cd799439011/eissue/status \
  -H "Authorization: Bearer token"

# Should return:
# {
#   "status": "DRAFT" or "ISSUED" or "FAILED",
#   "provider": "gib",
#   "external_id": "UUID-from-GIB",
#   "sent_at": "2026-02-02T10:00:00Z"
# }
```

---

## 🚀 Production Deployment

### 1. Pre-Deployment Validation

**Checklist:**
- [ ] All sandbox tests passing
- [ ] Production credentials obtained from GIB
- [ ] Certificates downloaded and validated
- [ ] Environment variables configured in production
- [ ] Backup/disaster recovery plan documented
- [ ] Monitoring & alerting configured (see below)

### 2. Environment Setup

**Option A: Kubernetes Secrets**

```bash
# Create secret for GIB credentials
kubectl create secret generic gib-credentials \
  --from-literal=GIB_USERNAME=prod_username \
  --from-literal=GIB_PASSWORD=prod_password \
  --from-literal=GIB_VKN_TCKN=12345678900 \
  --from-file=GIB_CERT_PATH=./certificate.cer \
  --from-file=GIB_PRIVATE_KEY_PATH=./private.key

# Reference in deployment
spec:
  containers:
  - name: api
    env:
    - name: GIB_ENVIRONMENT
      value: "prod"
    - name: GIB_BASE_URL
      value: "https://efatura.gib.gov.tr/v1/soap"
    - name: GIB_USERNAME
      valueFrom:
        secretKeyRef:
          name: gib-credentials
          key: GIB_USERNAME
    - name: GIB_PASSWORD
      valueFrom:
        secretKeyRef:
          name: gib-credentials
          key: GIB_PASSWORD
    # ... etc
    volumeMounts:
    - name: gib-certs
      mountPath: /secure/certs
      readOnly: true
  volumes:
  - name: gib-certs
    secret:
      secretName: gib-credentials
      defaultMode: 0400
```

**Option B: AWS Secrets Manager**

```go
// Load from AWS Secrets Manager (in production)
import "github.com/aws/aws-secretsmanager-caching-go/cache"

secretCache := cache.New(
    secretsmanageriface.NewSecretsManagerClient(sess),
    cache.CacheConfig{},
)

gibPassword := secretCache.GetSecretString("gib/prod/password")
```

### 3. Deployment Steps

1. **Deploy new code:**
   ```bash
   # Ensure e-fatura service methods are in production version
   git tag v2.0.0-eifatura-prod
   git push origin v2.0.0-eifatura-prod
   ```

2. **Configure production environment:**
   ```bash
   # Update production .env or Kubernetes secrets
   # Run configuration validation
   curl http://api:8080/health/config | grep gib
   ```

3. **Perform production smoke tests:**
   ```bash
   # Test with real GIB production endpoint
   # Send test invoice (low-value transaction)
   # Verify status query returns correct response
   # Verify cancellation works
   ```

4. **Monitor for issues:**
   ```bash
   # Watch application logs
   tail -f /var/log/api/eifatura.log | grep ERROR
   
   # Check metrics
   # - e-invoice send success rate (target: >99%)
   # - GIB API response time (target: <5s)
   # - Error rates by type (validation, network, auth)
   ```

---

## 🔍 Monitoring & Alerts

### Metrics to Track

```go
// Key metrics to export to monitoring system
metrics := []string{
  "einvoice_submit_total",        // Total invoices submitted
  "einvoice_submit_success",      // Successful submissions
  "einvoice_submit_failure",      // Failed submissions
  "einvoice_submit_duration_ms",  // Response time
  "einvoice_cancel_total",        // Total cancellations
  "einvoice_status_query_total",  // Status queries
  "gib_api_errors_by_type",       // Error breakdown (validation, auth, network)
}
```

### Alert Rules

**Critical (Page on-call):**
```
- e-invoice send success rate < 95% for 5 minutes
- GIB API connection errors > 10% for 2 minutes
- Certificate validation failures
```

**Warning (Create incident):**
```
- e-invoice send success rate < 98%
- Average response time > 3 seconds
- Daily failed invoices > 5
```

### Dashboard Example

```
📊 e-Fatura Health Dashboard

┌─────────────────────────────────────┐
│ Last 24 Hours                       │
├─────────────────────────────────────┤
│ Invoices Submitted:    1,234        │
│ Success Rate:          99.2% ✅     │
│ Avg Response Time:     2.3s         │
│ Failed Invoices:       9            │
│ Cancelled Invoices:    3            │
└─────────────────────────────────────┘

Failures by Type:
├─ Validation Error:    4 (44%)
├─ Network Error:       3 (33%)
├─ Auth Error:          1 (11%)
└─ Other:              1 (11%)
```

---

## 🛠️ Troubleshooting

### Common Issues

**Issue 1: "Invalid Credentials"**
```
Error: gib.auth_failed
Solution:
- Verify GIB_USERNAME and GIB_PASSWORD are correct
- Check that credentials are for PRODUCTION, not sandbox
- Confirm account is active with GIB
- Retry after 5 minutes (may be rate limited)
```

**Issue 2: "Certificate Not Valid"**
```
Error: gib.signature_failed
Solution:
- Verify certificate dates: openssl x509 -in cert.cer -noout -dates
- Ensure certificate is from authorized GIB provider
- Check private key matches certificate
- Verify file paths are correct
```

**Issue 3: "Invoice Already Submitted"**
```
Error: gib.duplicate_invoice
Solution:
- Idempotency: check `external_id` before resubmitting
- If invoice exists, query status instead of resubmit
- Store successful external_id in invoice.eissue.external_id
```

**Issue 4: "Network Timeout"**
```
Error: context deadline exceeded
Solution:
- Increase timeout: DEFAULT is 30s, can set 60s for prod
- Check GIB service status: https://efatura.gib.gov.tr/status
- Verify network connectivity: curl https://efatura.gib.gov.tr/
- Check firewall rules for outbound HTTPS to GIB
```

### Production Runbook

**When e-invoice submission fails:**

1. Check GIB service status dashboard
2. Verify network connectivity to GIB endpoint
3. Check certificate expiration
4. Review error logs: `grep ERROR /var/log/api/eifatura.log`
5. For transient errors: automatic retry via DLQ after 1 hour
6. For persistent errors: escalate to GIB support with request ID

**When to contact GIB Support:**
- Certificate errors that don't resolve after renewal
- Account locked or suspended
- Persistent authentication failures
- API endpoint unavailability > 15 minutes

GIB Support: https://efatura.gib.gov.tr/support

---

## 📝 Operational Procedures

### Daily Tasks
- [ ] Monitor e-invoice dashboard (8:00 AM)
- [ ] Check for failed invoices (review/retry as needed)
- [ ] Verify certificate expiration (alerts if < 30 days)

### Weekly Tasks
- [ ] Review error logs for patterns
- [ ] Reconcile invoice counts with GIB reports
- [ ] Test status query & cancellation operations

### Monthly Tasks
- [ ] Review GIB compliance requirements
- [ ] Update operational runbook based on incidents
- [ ] Capacity planning (invoice volume trends)

### Quarterly Tasks
- [ ] Certificate renewal (before expiration)
- [ ] Security audit of credential storage
- [ ] Disaster recovery drill

---

## 🔐 Security Best Practices

1. **Credential Management:**
   - Never commit credentials to Git
   - Use secrets manager (Kubernetes, AWS, HashiCorp Vault)
   - Rotate credentials quarterly
   - Use strong passwords (32+ characters)

2. **Certificate Management:**
   - Use hardware security modules (HSM) for private keys
   - Set file permissions: `chmod 600` on private key
   - Audit access to certificate files
   - Maintain certificate chain for validation

3. **API Communication:**
   - Use HTTPS only (TLS 1.2+)
   - Validate certificate chain
   - Implement request signing
   - Rate limit requests per GIB guidelines

4. **Audit & Compliance:**
   - Log all invoice operations (audit trail)
   - PII masking in logs
   - Compliance with GIB data retention policies
   - Regular security assessments

---

## 📞 Support Contacts

| Role | Contact | Availability |
|------|---------|--------------|
| **GIB Technical Support** | https://efatura.gib.gov.tr/support | 08:00-17:00 TR Time, Mon-Fri |
| **GIB Helpline** | +90 (312) 204-90-00 ext. 1234 | Business hours |
| **Our e-Fatura Team** | See on-call schedule | 24/7 for critical issues |
| **E-Fatura Lead** | See organization chart | Business hours |

---

## 🎯 Success Criteria

**Production is stable when:**
- ✅ E-invoice success rate > 99% sustained for 48+ hours
- ✅ All invoices can be queried for status from GIB
- ✅ Cancellations process without errors
- ✅ Alerts are configured and working
- ✅ Runbook covers all known failure scenarios
- ✅ Team trained on operational procedures

**Target launch date:** Week 1 post-launch (after GIB prod credentials received)

---

**Document Version:** 1.0 (Production Ready)  
**Last Updated:** February 2, 2026  
**Next Review:** After first production month
