import matplotlib.pyplot as plt

def plot_team_stats(stats, team_name):
    categories = ['Wins', 'Losses', 'Draws']
    values = [stats.get('wins', 0), stats.get('losses', 0), stats.get('draws', 0)]

    plt.bar(categories, values, color=['green', 'red', 'gray'])
    plt.title(f"Performance of {team_name}")
    plt.ylabel("Number of Matches")
    plt.show()

def plot_score_trend(matches, team):
    dates = list(range(len(matches)))
    scores = []

    for m in matches:
        if m["team_a"] == team:
            scores.append(m["score_a"])
        elif m["team_b"] == team:
            scores.append(m["score_b"])

    plt.plot(dates, scores, marker='o')
    plt.title(f"{team} - Score Trend")
    plt.xlabel("Match Number")
    plt.ylabel("Goals Scored")
    plt.show()
