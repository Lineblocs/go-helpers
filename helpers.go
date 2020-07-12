package lineblocs 

import (
    "net/http"
	"os"
	"time"
	"strconv"
	"net"
	"strings"
	"context"
	//"errors"
	"mime/multipart"
	"reflect"
	"fmt"
	"database/sql"
	"regexp"
	_ "github.com/go-sql-driver/mysql"
	guuid "github.com/google/uuid"
	libphonenumber "github.com/ttacon/libphonenumber"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/mailgun/mailgun-go/v4"
)

type Call struct {
  From string `json:"from"`
  To string `json:"to"`
  Status string `json:"status"`
  Direction string `json:"direction"`
  Duration string `json:"duration"`
  UserId int `json:"user_id"`
  WorkspaceId int  `json:"workspace_id"`
  APIId string `json:"api_id"`
}
type CallUpdateReq struct {
  CallId int `json:"call_id"`
  Status string `json:"status"`
}
type RecordingTranscriptionReq struct {
	RecordingId int `json:"recording_id"`
  Ready bool `json:"ready"`
  Text string `json:"text"`
}
type Conference struct {
  Name string `json:"name"`
  WorkspaceId int `json:"workspace_id"`
  APIId string `json:"api_id"`
}
type DebitCreateReq struct {
  UserId int `json:"user_id"`
  WorkspaceId int `json:"workspace_id"`
  ModuleId int `json:"module_id"`

  Source string `json:"source"`
  Number string `json:"number"`
  Type string `json:"type"`
  Seconds float64 `json:"seconds"`
}


type CallRate struct {
	CallRate float64
}


type DebitAPIParams struct {
	Length int `json:"length"`	
	RecordingLength float64 `json:"recording_length"`	
}
type DebitAPICreateReq struct {
  UserId int `json:"user_id"`
  WorkspaceId int `json:"workspace_id"`
  Type string `json:"type"`
  Source string `json:"source"`
  Params DebitAPIParams `json:"params"`
}

type LogCreateReq struct {
  UserId int `json:"user_id"`
  WorkspaceId int `json:"workspace_id"`
  Title string `json:"title"`
  Report string `json:"report"`
  FlowId int `json:"flow_id"`
  Level *string `json:"report"`
  From *string `json:"from"`
  To *string `json:"to"`
}
type LogSimpleCreateReq struct {
  Type string `json:"type"`
  Level *string `json:"level"`
}
type Fax struct {
  UserId int `json:"user_id"`
  WorkspaceId int `json:"workspace_id"`
  CallId int `json:"call_id"`
  Uri string `json:"uri"`
  APIId string `json:"api_id"`
}

type Recording struct {
  Id int `json:"id"`
  UserId int `json:"user_id"`
  CallId int `json:"call_id"`
  Size int `json:"size"`
  WorkspaceId int `json:"workspace_id"`
  APIId string `json:"api_id"`
  Tags *[]string `json:"tags"`
	TranscriptionReady bool `json:"transcription_ready"`
	TranscriptionText string `json:"transcription_text"`
}

type VerifyNumber struct {
	Valid bool `json:"valid"`
}




type LogRoutine struct {
  UserId int
  WorkspaceId int
  Title string
  Report string
  FlowId int
  Level string
  From string
  To string
}
type User struct {
  Id int
  Username string
  FirstName string
  LastName string
  Email string
}

type Workspace struct {
  Id int `json:"id"`
  CreatorId int `json:"creator_id"`
  Name string `json:"name"`
  BYOEnabled bool `json:"byo_enabled"`
  IPWhitelistDisabled bool `json:"ip_whitelist_disabled"`
  OutboundMacroId int `json:"outbound_macro_id"`
  Plan string `json:"plan"`
}

type WorkspaceParam struct {
	Key string `json:"key"`
	Value string `json:"value"`
}
type WorkspaceCreatorFullInfo struct {
  Id int `json:"id"`
	Workspace *Workspace `json:"workspace"`
	WorkspaceName string `json:"workspace_name"`
	WorkspaceDomain string `json:"workspace_domain"`
	WorkspaceId int `json:"workspace_id"`
	WorkspaceParams *[]WorkspaceParam `json:"workspace_params"`
  	OutboundMacroId int `json:"outbound_macro_id"`
}
type WorkspaceDIDInfo struct {
  WorkspaceId int `json:"workspace_id"`
  Number string `json:"number"`
  FlowJSON string `json:"flow_json"`
  WorkspaceName string `json:"workspace_name"`
  Name string `json:"name"`
  Plan string `json:"plan"`
  BYOEnabled bool `json:"byo_enabled"`
  IPWhitelistDisabled bool `json:"ip_whitelist_disabled"`
  OutboundMacroId int `json:"outbound_macro_id"`
  CreatorId int `json:"creator_id"`
  APIToken string `json:"api_token"`
  APISecret string `json:"api_secret"`
  WorkspaceParams *[]WorkspaceParam `json:"workspace_params"`
}
type WorkspacePSTNInfo struct {
  IPAddr string `json:"ip_addr"`
  DID string `json:"did"`
}
type CallerIDInfo struct {
  CallerID string `json:"caller_id"`
}
type ExtensionFlowInfo struct {
  CallerID string `json:"caller_id"`
  WorkspaceId int `json:"workspace_id"`
  FlowJSON string `json:"flow_json"`
  Username string `json:"username"`
  Name string `json:"name"`
  WorkspaceName  string `json:"workspace_name"`
  Plan string `json:"plan"`
  CreatorId int `json:"creator_id"`
  Id int `json:"id"`
  APIToken string `json:"api_token"`
  APISecret string `json:"api_secret"`
  WorkspaceParams *[]WorkspaceParam `json:"workspace_params"`
  FreeTrialStatus string `json:"workspace_params"`
}

type CodeFlowInfo struct {
  WorkspaceId int `json:"workspace_id"`
  Code string `json:"code"`
  FlowJSON string `json:"flow_json"`
  Name string `json:"name"`
  WorkspaceName  string `json:"workspace_name"`
  Plan string `json:"plan"`
  CreatorId int `json:"creator_id"`
  Id int `json:"id"`
  APIToken string `json:"api_token"`
  APISecret string `json:"api_secret"`
  FreeTrialStatus string `json:"workspace_params"`
  FoundCode bool `json:"found_code"`
}


type MacroFunction struct {
	Title string `json:"title"`
	Code string `json:"code"`
	CompiledCode string `json:"compiled_code"`
}
type MediaServer struct {
	IpAddress string `json:"ip_address"`
	PrivateIpAddress string `json:"private_ip_address"`
}
type EmailInfo struct {
	Message string `json:"message"`
}

type GlobalSettings struct {
  ValidateCallerId bool
}

type ServicePlan struct {

}

var db* sql.DB;
var settings *GlobalSettings;
func CreateDBConn() (*sql.DB, error) {
	if db != nil {
		return db, nil
	}
	db, err := sql.Open("mysql", `lineblocs:&!UER~7$Z>fx3S3J@tcp(lineblocs.ckehyurhpc6m.ca-central-1.rds.amazonaws.com:3306)/lineblocs?parseTime=true`)
	//db, err = sql.Open("mysql", "root:mysql@tcp(127.0.0.1:3306)/lineblocs?parseTime=true") //add parse time
	if err != nil {
		panic("Could not connect to MySQL");
		return nil, err
	}
	  db.SetMaxOpenConns(10)
	  return db, nil
}

func CreateAPIID(prefix string) string {
	id := guuid.New()
	return prefix + "-" + id.String()
}
func LookupBestCallRate(number string, typeRate string) *CallRate {
	return &CallRate{ CallRate: 9.99 };
}
func HandleInternalErr(msg string, err error, w http.ResponseWriter) {
	fmt.Printf(msg)
	fmt.Println(err)
	w.WriteHeader(http.StatusInternalServerError)
}
func CalculateTTSCosts(length int) float64 {
	var result float64 = float64(length) * .000005
	return result
}
func CalculateSTTCosts(recordingLength float64) float64 {
	// Google cloud bills .006 per 15 seconds
	billable := recordingLength / 15
	var result float64 = 0.006 * billable
	return result
}
func GetUserFromDB(id int) (*User, error) {
	var userId int
	var username string
	var fname string
	var lname string
	var email string
	fmt.Printf("looking up user %d\r\n", id)
	row := db.QueryRow(`SELECT id, username, first_name, last_name, email FROM users WHERE id=?`, id)

	err := row.Scan(&userId, &username, &fname, &lname,  &email)
	if ( err == sql.ErrNoRows ) {  //create conference
		return nil, err
	}
	if ( err != nil ) {  //another error
		return nil, err
	}

	return &User{Id: userId, Username: username, FirstName: fname, LastName: lname, Email: email}, nil
}
func GetWorkspaceFromDB(id int) (*Workspace, error) {
	var workspaceId int
	var name string
	var creatorId int
	var outboundMacroId sql.NullInt64
	var plan string
	row := db.QueryRow(`SELECT id, name, creator_id, outbound_macro_id, plan FROM workspaces WHERE id=?`, id)

	err := row.Scan(&workspaceId, &name, &creatorId, &outboundMacroId, &plan)
	if ( err == sql.ErrNoRows ) {  //create conference
		return nil, err
	}
	if ( err != nil ) {  //another error
		return nil, err
	}
    if reflect.TypeOf(outboundMacroId) == nil {
		return &Workspace{Id: workspaceId, Name: name, CreatorId: creatorId, Plan: plan}, nil
	}
	return &Workspace{Id: workspaceId, Name: name, CreatorId: creatorId, OutboundMacroId: int(outboundMacroId.Int64), Plan: plan}, nil
}

func GetRecordingSpace(id int) (int, error) {
	var bytes int
	row := db.QueryRow(`SELECT SUM(size) FROM recordings WHERE workspace_id=?`, id)

	err := row.Scan(&bytes)
	if ( err == sql.ErrNoRows ) {  //create conference
		return 0, err
	}
	if ( err != nil ) {  //another error
		return 0, err
	}
	return bytes, nil
}
func GetFaxCount(id int) (*int, error) {
	var count int
	row := db.QueryRow(`SELECT COUNT(*) FROM faxes WHERE workspace_id=?`, id)

	err := row.Scan(&count)
	if ( err == sql.ErrNoRows ) {  //create conference
		return nil, err
	}
	if ( err != nil ) {  //another error
		return nil, err
	}
	return &count, nil
}
func GetWorkspaceByDomain(domain string) (*Workspace, error) {
	var workspaceId int
	var name string
	var byo bool
	var ipWhitelist bool
	var creatorId int
	s := strings.Split(domain, ".")
	workspaceName := s[0]
	row := db.QueryRow("SELECT id, creator_id, name, byo_enabled, ip_whitelist_disabled FROM workspaces WHERE name=?", workspaceName)

	err := row.Scan(&workspaceId, &creatorId, &name, &byo, &ipWhitelist)
	if ( err == sql.ErrNoRows ) {  //create conference
		return nil, err
	}
	return &Workspace{Id: workspaceId, CreatorId: creatorId, Name: name, BYOEnabled: byo, IPWhitelistDisabled: ipWhitelist}, nil
}

func GetWorkspaceParams(workspaceId int) (*[]WorkspaceParam, error) {
	// Execute the query
	results, err := db.Query("SELECT `key`, `value` FROM workspace_params WHERE `workspace_id` = ?", workspaceId)
    if err != nil {
		return nil, err;
	}
  defer results.Close()
	params := []WorkspaceParam{};

    for results.Next() {
		param := WorkspaceParam{};
        // for each row, scan the result into our tag composite object
        err = results.Scan(&param.Key, &param.Value)
        if err != nil {
			return nil, err
		}
		params = append(params, param)
	}
	return &params, nil;
}

func GetUserByDomain(domain string) (*WorkspaceCreatorFullInfo, error) {
	workspace, err := GetWorkspaceByDomain( domain )
	if err != nil {
		return nil, err
	}

	// Execute the query
	params, err  := GetWorkspaceParams(workspace.Id)
    if err != nil {
		return nil, err;
	}
	full := &WorkspaceCreatorFullInfo{ 
    Id: workspace.CreatorId,
    Workspace: workspace, 
		WorkspaceParams: params,
		WorkspaceName: workspace.Name,
		WorkspaceDomain: fmt.Sprintf("%s.lineblocs.com", workspace.Name),
		WorkspaceId: workspace.Id,
		OutboundMacroId: workspace.OutboundMacroId	};

	return full, nil
}

func GetRecordingFromDB(id int) (*Recording, error) {
	var apiId string
	var ready int
	var size int
	var text string
	row := db.QueryRow("SELECT api_id, transcription_ready, transcription_text, size FROM recordings WHERE id=?", id)

	err := row.Scan(&apiId, &ready, &text, &size)
	if ( err == sql.ErrNoRows ) {  //create conference
		return nil, err
	}
	if ready == 1 {
		return &Recording{APIId: apiId, Id: id, TranscriptionReady: true, TranscriptionText: text, Size: size}, nil
	}
	return &Recording{APIId: apiId, Id: id, Size: size}, nil
}
//todo move to microservice
func GetPlanRecordingLimit(workspace* Workspace) (int, error) {
	if workspace.Plan == "pay-as-you-go" {
		return 1024, nil
	} else if workspace.Plan == "starter" {
		return 1024*2, nil
	} else if workspace.Plan == "pro" {
		return 1024*32, nil
	} else if workspace.Plan == "starter" {
		return 1024*128, nil
	}
	return 0, nil
}
//todo move to microservice
func GetPlanFaxLimit(workspace* Workspace) (*int, error) {
	var res* int
	if workspace.Plan == "pay-as-you-go" {
		*res = 100
	} else if workspace.Plan == "starter" {
		*res = 100
	} else if workspace.Plan == "pro" {
		res =  nil
	} else if workspace.Plan == "starter" {
		res = nil
	}
	return res, nil
}
func SendLogRoutineEmail(log* LogRoutine, user* User, workspace* Workspace) error {
	mg := mailgun.NewMailgun(os.Getenv("MAILGUN_DOMAIN"),os.Getenv("MAILGUN_API_KEY"))
	m := mg.NewMessage(
		"Lineblocs <monitor@lineblocs.com>",
		"Debug Monitor",
		"Debug Monitor",
		user.Email)
	m.AddCC("contact@lineblocs.com")
	//m.AddBCC("bar@example.com")


	body := `<html>
<head></head>
<body>
	<h1>Lineblocs Monitor Report</h1>
	<h5>` + log.Title + `</h5>
	<p>` + log.Report + `</p>
</body>
</html>`;

	m.SetHtml(body)
	//m.AddAttachment("files/test.jpg")
	//m.AddAttachment("files/test.txt")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, _, err := mg.Send(ctx, m)
	if err != nil {
		return err
	}
	return nil
}

func StartLogRoutine(log* LogRoutine) (*string, error) {
	var user* User;
	var workspace* Workspace;

    user, err := GetUserFromDB(log.UserId)
	if err != nil {
		fmt.Printf("could not get user..")
		return nil, err
	}

	workspace, err = GetWorkspaceFromDB(log.WorkspaceId)
	if err != nil {
		fmt.Printf("could not get workspace..")
		return nil, err
	}
	now := time.Now()
	apiId := CreateAPIID("log")
	stmt, err := db.Prepare("INSERT INTO debugger_logs (`from`, `to`, `title`, `report`, `workspace_id`, `level`, `api_id`, `created_at`, `updated_at`) VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ? )")

	if err != nil {
		fmt.Printf("could not prepare query..")
		return nil, err
	}

	defer stmt.Close()
	res, err := stmt.Exec(log.From, log.To, log.Title, log.Report, workspace.Id, log.Level, apiId, now, now)
	if err != nil {
		fmt.Printf("could not execute query..")
		return nil, err
	}

	logId, err := res.LastInsertId()
	if err != nil {
		fmt.Printf("could not get insert id..")
		return nil, err
	}
	logIdStr := strconv.FormatInt(logId, 10)

	go SendLogRoutineEmail(log, user, workspace)

	return &logIdStr, err
}
func CheckRouteMatches(from string, to string, prefix string, prepend string, match string) (bool, error) {
	full := prefix + match
	valid, err := regexp.MatchString(full, to)
	if err != nil {
		return false, err
	}
	return valid, err
}
func ShouldUseProviderNext(name string, ipPrivate string) (bool, error) {
	return true, nil
}
func CheckCIDRMatch(sourceIp string, fullIp string) (bool, error) {
	_, net1, err :=  net.ParseCIDR(sourceIp + "/32")
	if err != nil {
		return false, err
	}
	_, net2, err :=  net.ParseCIDR(fullIp)
	if err != nil {
		return false, err
	}

	return net2.Contains(net1.IP), nil
}
func CheckPSTNIPWhitelist(did string, sourceIp string) (bool, error) {
	results, err := db.Query(`SELECT 
	sip_providers_whitelist_ips.ip_address, 
	sip_providers_whitelist_ips.ip_address_range
	FROM sip_providers_whitelist_ips
	INNER JOIN sip_providers ON sip_providers.id = sip_providers_whitelist_ips.provider_id
	INNER JOIN did_numbers ON did_numbers.workspace_id = sip_providers_whitelist_ips.workspace_id
	WHERE did_numbers.api_number = ?
	`, did)
    if err != nil {
		return false, err
	}
  defer results.Close()
    for results.Next() {
		var ipAddr string
		var ipAddrRange string
		err = results.Scan(&ipAddr, &ipAddrRange)
		if err != nil {
		  return false, err

		}
		fullIp := ipAddr + ipAddrRange
		match, err := CheckCIDRMatch(sourceIp, fullIp) 
		if err != nil {
		  fmt.Printf("error matching CIDR source %s, full %s\r\n", sourceIp, fullIp)
		  continue
		}
		if match {
			return true, nil
		}
	}
	return false, nil
}
func CheckBYOPSTNIPWhitelist(did string, sourceIp string) (bool, error) {
	results, err := db.Query(`SELECT 
	byo_carriers_ips.ip,
	byo_carriers_ips.range
	FROM byo_carriers_ips
	INNER JOIN byo_carriers ON byo_carriers.id = byo_carriers_ips.carrier_id
	INNER JOIN byo_did_numbers ON byo_did_numbers.workspace_id = byo_carriers.workspace_id
	WHERE byo_did_numbers.number = ?
	`, did)
    if err != nil {
		return false, err
	}
  defer results.Close()
    for results.Next() {
		var ipAddr string
		var ipAddrRange string
		err = results.Scan(&ipAddr, &ipAddrRange)
		if err != nil {
			return false, err
		}
		fullIp := ipAddr + ipAddrRange
		match, err := CheckCIDRMatch(sourceIp, fullIp) 
		if err != nil {
		  fmt.Printf("error matching CIDR source %s, full %s\r\n", sourceIp, fullIp)
		  continue
		}
		if match {
			return true, nil
		}
	}
	return false, nil
}

func FinishValidation(number string, didWorkspaceId string) (bool,error) {
	num, err := libphonenumber.Parse(number, "US")
	if err != nil {
		return false, err
	}
	formattedNum := libphonenumber.Format(num, libphonenumber.E164)
	row := db.QueryRow("SELECT id FROM `blocked_numbers` WHERE `workspace_id` = ? AND `number` = ?", didWorkspaceId, formattedNum)
	var id string
	err = row.Scan(&id);
	if ( err == sql.ErrNoRows ) {  //create conference
		return true,nil
	}
	if err != nil {
		return false,err
	}
	return false,nil
}
func CheckFreeTrialStatus(plan string, started time.Time) string {
	if plan  == "trial" {
		now := time.Now()
		//make configurable
		expireDays := 10
		expireHours := expireDays * 24
		started.Add(time.Hour * time.Duration(expireHours))
		if started.After( now ) {
			return "expired";
		}
		return "pending-expiry";
	}
	return "not-applicable";
}
func CheckIsMakingOutboundCallFirstTime(call Call) {
	var id string
	row := db.QueryRow("SELECT id FROM `calls` WHERE `workspace_id` = ? AND `from` LIKE '?%s' AND `direction = 'outbound'", call.WorkspaceId, call.From, call.Direction)
	err := row.Scan(&id);
	if ( err != sql.ErrNoRows ) {  //create conference
		// all ok
		return
	}
	//send notification
	user, err := GetUserFromDB(call.UserId)
	if err != nil {
		panic(err)
	}
	body := `A call was made to ` + call.To + ` for the first time on your account.`;
	SendEmail(user, "First call to destination country", body)
}
func SendEmail(user *User, subject string, body string) {
}
func SomeLoadBalancingLogic() (*MediaServer,error) {
	results, err := db.Query("SELECT ip_address,private_ip_address FROM media_servers");
    if err != nil {
		return nil,err
	}
  defer results.Close()
    for results.Next() {
		value := MediaServer{};
		err = results.Scan(&value.IpAddress,&value.PrivateIpAddress);
		if err != nil {
			return nil,err
		}
		return &value,nil
	}
	return nil,nil
}
func DoVerifyCaller(workspaceId int, number string) (bool, error) {
	var workspace* Workspace;

  if !settings.ValidateCallerId { 
    return true, nil
  }

	workspace, err := GetWorkspaceFromDB(workspaceId)
	if err != nil {
		return false, err
	}

	num, err := libphonenumber.Parse(number, "US")
	if err != nil {
		return false, err
	}
	formattedNum := libphonenumber.Format(num, libphonenumber.E164)
	fmt.Printf("looking up number %s\r\n", formattedNum)
	fmt.Printf("domain isr %s\r\n", workspace.Name)
	var id string
	row := db.QueryRow("SELECT id FROM `did_numbers` WHERE `number` = ? AND `workspace_id` = ?", formattedNum, workspace.Id)
	err = row.Scan(&id);
	if ( err != sql.ErrNoRows ) {  //create conference
		return true, nil
	}
	return false, nil
}

func GetQueryVariable(r *http.Request, key string) *string {
	vals := r.URL.Query() // Returns a url.Values, which is a map[string][]string
	results, ok := vals[key] // Note type, not ID. ID wasn't specified anywhere.
	var value *string
	if ok {
		if len(results) >= 1 {
			value = &results[0] // The first `?type=model`
		}
	}
	return value
}
func UploadS3(folder string, name string, file multipart.File) (error) {
	bucket := "lineblocs"
	key := folder + "/" + name
	// The session the S3 Uploader will use
	session, err := session.NewSession(&aws.Config{
		Region: aws.String("ca-central-1")})
	if err != nil {
		return fmt.Errorf("S3 session err: %s", err)
	}

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(session)

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}
	fmt.Printf("file uploaded to, %s\n", aws.StringValue(&result.Location))
	return nil
}

func CreateS3URL(folder string, id string) string {
	return "https://lineblocs.s3.ca-central-1.amazonaws.com/" + folder + "/" + id
}

func NoContent(w http.ResponseWriter, r *http.Request) {
  // Set up any headers you want here.
  w.WriteHeader(http.StatusNoContent) // send the headers with a 204 response code.
}
func ToCents(dollars float64) int {
	result := dollars * 100
	return int( result )
}

func GetServicePlans() ([]ServicePlan, error) {
	plans := []ServicePlan{};
	plans = append(plans, ServicePlan{})
	return plans, nil
}
