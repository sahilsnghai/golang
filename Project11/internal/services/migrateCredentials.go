package services

import (
	"fmt"
	"log"
	"strings"
)

func (s *Storage) MigrateCreds(data map[string]interface{}) (any, error) {
	log.Println("Migrating Credentials")
	domainIDStr := fmt.Sprintf("%0.f", data["domainId"])

	rows, err := s.QueryWithTiming(`SELECT
				host_ip as host,
				host_port as port,
				user_name as user,
				backup_host as backup_server_node,
				user_password as "password",
				database_name as "database",
				database_driver as driver
			from
				sch_connection_detail scd
			join sch_schema ss on
				scd.sch_connection_detail_id = ss.sch_connection_id
			where
				ss.sch_schema_id = ?;`, data["domainId"])
	if err != nil {
		return nil, fmt.Errorf("unable to fetch %w", err)
	}

	responseCredentials, err := GetFromRows(rows)

	if err != nil {
		return nil, fmt.Errorf("unable to fetch or convert creds. %w", err)
	}
	resp := make(map[string]interface{})

	if len(responseCredentials) > 0 {
		resp = responseCredentials[0]

		if val, ok := resp["backup_server_node"].(string); ok && val != "" {
			resp["backup_server_node"] = strings.Split(val, ",")
		} else {
			resp["backup_server_node"] = []string{}
		}
	}
	go s.hPublish(domainIDStr, "credentials", resp)
	return resp, nil
}
