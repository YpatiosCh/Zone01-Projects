package graphql

const (
	// normal with arguments
	USER_INFO = `
		query ($userId: Int!) {
			user(where: { id: { _eq: $userId }}) {
				id
				login
				campus
				createdAt
				attrs
			}
		}
	`
	// normal with arguments
	USER_XP = `
		query ($userId: Int!) {
  			transaction_aggregate(
    			where: {
      				userId: { _eq: $userId },
      				type: { _eq: "xp" },
      				path: { _regex: "^/athens/div-01/[^/]+$" },
      				campus: { _eq: "athens" }
    			},
    			distinct_on: path
  			) {
    		   aggregate {
      				sum { amount }
    		   }		
  			}
		}
	`
	// normal with arguments
	USER_PROGRESS = `
		query ($userId: Int!) {
  			transaction(
    			where: {
      				userId: { _eq: $userId },
      				type: { _eq: "xp" },
      				path: { _regex: "^/athens/div-01/[^/]+$" },
      				campus: { _eq: "athens" }
    			},
    			distinct_on: path
  			) {
				path
				amount
				createdAt
  			}
		}	
	`
	// normal with arguments
	USER_GROUP_IDS = `
	query ($userId: Int!) {
  		group_user(
    		where: {
      			userId: { _eq: $userId }
				accepted: { _eq: true }
      			group: { path: { _regex: "^/athens/div-01/[^/]+$" } }
    		}
  		) {
    		groupId
  		}
	}
	`

	USER_PROJECTS = `
	query ($userId: Int!) {
  		transaction(
    		where: {
      			userId: { _eq: $userId },
      			type: { _eq: "xp" },
      			path: { _regex: "^/athens/div-01/[^/]+$" },
      			campus: { _eq: "athens" }
    		}
        ) {
    		path
  		}			
  	}
	`

	// nested query with arguments
	USER_COLLABORATORS = `
	query GetCollaborators($userId: Int!, $groupIds: [Int!]) {
  		group_user(
    		where: {
      			groupId: { _in: $groupIds }
				accepted: { _eq: true }
      			userId: { _neq: $userId }
    		}
  		) {
			user {
				login
			}
		}
	}
	`
)
