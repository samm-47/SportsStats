# client.py
import requests
import csv
import json

BASE_URL = 'http://localhost:8080'

def add_match():
    team_a = input("Team A: ")
    score_a = int(input("Score A: "))
    team_b = input("Team B: ")
    score_b = int(input("Score B: "))
    data = {
        "team_a": team_a,
        "score_a": score_a,
        "team_b": team_b,
        "score_b": score_b
    }
    res = requests.post(f'{BASE_URL}/match', json=data)
    print(res.text)

def get_team_stats():
    team = input("Enter team name: ")
    res = requests.get(f'{BASE_URL}/team/{team}')
    print(res.json())

def batch_upload_csv(file_path):
    with open(file_path, newline='') as csvfile:
        reader = csv.DictReader(csvfile)
        for row in reader:
            data = {
                "team_a": row['team_a'],
                "score_a": int(row['score_a']),
                "team_b": row['team_b'],
                "score_b": int(row['score_b'])
            }
            res = requests.post(f'{BASE_URL}/match', json=data)
            print(res.text)

def batch_upload_json(file_path):
    with open(file_path) as f:
        matches = json.load(f)
        for match in matches:
            data = {
                "team_a": match['team_a'],
                "score_a": int(match['score_a']),
                "team_b": match['team_b'],
                "score_b": int(match['score_b'])
            }
            res = requests.post(f'{BASE_URL}/match', json=data)
            print(res.text)

def predict_match():
    team_a = input("Enter Team A: ")
    team_b = input("Enter Team B: ")
    res = requests.get(f'{BASE_URL}/predict/{team_a}/{team_b}')
    if res.status_code == 200:
        data = res.json()
        print(f"Probability {team_a} wins: {data['prob_team_a_win']:.2%}")
        print(f"Probability {team_b} wins: {data['prob_team_b_win']:.2%}")
    else:
        print("Error getting prediction", res.text)

def main():
    while True:
        print("\n1. Add match")
        print("2. Get team stats")
        print("3. Upload CSV")
        print("4. Upload JSON")
        print("5. Predict match outcome")
        print("6. Exit")
        choice = input("Choose an option: ")
        if choice == '1':
            add_match()
        elif choice == '2':
            get_team_stats()
        elif choice == '3':
            path = input("CSV file path: ")
            batch_upload_csv(path)
        elif choice == '4':
            path = input("JSON file path: ")
            batch_upload_json(path)
        elif choice == '5':
            predict_match()
        elif choice == '6':
            break
        else:
            print("Invalid choice")

if __name__ == '__main__':
    main()
