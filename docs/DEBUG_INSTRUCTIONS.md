# API Debugging Instructions

## Current Status
✅ **Fixed**: API port issue - killed duplicate Vite processes  
✅ **Added**: Comprehensive API debugging utility  
🔄 **In Progress**: Testing wallet display differences  

## How to Test

### 1. Open the Application
- Navigate to: http://localhost:3000
- The frontend should be running on port 3000 (not 3001)
- Backend should be accessible at http://localhost:8080

### 2. Use the API Tester
- Go to the Dashboard page
- Scroll down to find the "API Tester & Debugger" section
- Click "Run Tests" to test both wallets and categories endpoints
- Check the console for detailed API logs

### 3. Browser Console Debugging
Open browser console (F12) and use these commands:

```javascript
// Check for port issues
apiDebugger.detectPortIssues()

// Compare API calls between Dashboard and Wallets page
apiDebugger.compareComponents('Dashboard', 'WalletsPage')

// View all API logs
apiDebugger.getLogs()

// Clear debug logs
apiDebugger.clearLogs()
```

### 4. Compare Pages
1. **Dashboard Page**: Navigate to `/dashboard`
   - Check console for "🏠 Dashboard Render Debug" logs
   - Note how wallets are displayed in the "錢包概況" section
   
2. **Wallets Page**: Navigate to `/wallets` 
   - Check console for "🔄 Component render - wallets state" logs
   - Note how wallets are displayed in the main grid

### 5. Expected Results

#### Working Correctly:
- All API calls should go to `http://localhost:3000/api/v1/*` (not 3001)
- Wallets should load and display on both pages
- Categories endpoint should return 501 (Not Implemented) - this is expected

#### Issues to Look For:
- Wrong port usage (3001 instead of 3000)
- Different data structures between Dashboard and Wallets page
- Missing or failed API requests
- Console errors

## Key Files Modified

1. **`/src/utils/apiDebug.ts`** - New debugging utility
2. **`/src/services/api.ts`** - Enhanced with request/response logging
3. **`/src/services/walletService.ts`** - Added component tagging
4. **`/src/pages/Dashboard.tsx`** - Added debugging and component tagging
5. **`/src/pages/Wallets.tsx`** - Added component tagging
6. **`/src/components/ApiTester.tsx`** - New testing component

## Debugging Features

### API Request Tracking
- ✅ Full URL logging (detects wrong ports)
- ✅ Response time measurement
- ✅ Component-specific tagging
- ✅ Success/failure tracking
- ✅ Error details and status codes

### Console Utilities
- ✅ `apiDebugger.detectPortIssues()` - Port problem detection
- ✅ `apiDebugger.compareComponents()` - Compare API usage
- ✅ Color-coded console logs for easy reading
- ✅ Grouped console output for clarity

### Visual Testing
- ✅ In-app API tester component
- ✅ Real-time test results display
- ✅ One-click test execution

## Next Steps

1. Test the application with these tools
2. Identify any remaining issues with wallet display
3. Verify all API calls use correct port (3000)
4. Compare the exact differences between Dashboard and Wallets page
5. Fix any remaining inconsistencies

## Remove Debugging (Production)

When ready for production, remove:
- ApiTester component from Dashboard
- Console.log statements
- Debug instructions file